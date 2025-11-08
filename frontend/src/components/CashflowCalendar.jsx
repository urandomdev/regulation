import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { financialService } from '../services/financialService';
import { utils } from '../api/client';
import './CashflowCalendar.css';

const CashflowCalendar = () => {
  const navigate = useNavigate();
  const [weekDays, setWeekDays] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchWeekData = async () => {
      try {
        setLoading(true);

        // Get current week's date range (Monday to Sunday)
        const today = new Date();
        const dayOfWeek = today.getDay();
        const monday = new Date(today);
        monday.setDate(today.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1));
        const sunday = new Date(monday);
        sunday.setDate(monday.getDate() + 6);

        // Fetch transactions for the week
        const result = await financialService.getTransactions({
          start: utils.formatDate(monday),
          end: utils.formatDate(sunday),
          limit: 500,
        });

        // Group transactions by day
        const dayMap = {};
        result.data.forEach(t => {
          const date = new Date(t.date);
          const dayKey = utils.formatDate(date);

          if (!dayMap[dayKey]) {
            dayMap[dayKey] = { spend: 0, income: 0 };
          }

          if (t.isIncome) {
            dayMap[dayKey].income += t.amount;
          } else {
            dayMap[dayKey].spend += t.amount;
          }
        });

        // Create week days array
        const days = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'];
        const weekData = [];
        const todayStr = utils.formatDate(today);

        for (let i = 0; i < 7; i++) {
          const currentDay = new Date(monday);
          currentDay.setDate(monday.getDate() + i);
          const dayKey = utils.formatDate(currentDay);
          const dayData = dayMap[dayKey] || { spend: 0, income: 0 };

          weekData.push({
            day: days[currentDay.getDay()],
            date: currentDay.getDate(),
            spend: Math.round(dayData.spend),
            income: Math.round(dayData.income),
            isToday: dayKey === todayStr,
          });
        }

        setWeekDays(weekData);
      } catch (err) {
        console.error('Failed to fetch week data:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchWeekData();
  }, []);

  if (loading) {
    return (
      <div className="cashflow-calendar-card">
        <div className="card-header">
          <h2>This Week</h2>
        </div>
        <div className="loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="cashflow-calendar-card">
        <div className="card-header">
          <h2>This Week</h2>
        </div>
        <div className="error">Error: {error}</div>
      </div>
    );
  }

  const maxAmount = Math.max(...weekDays.map(d => Math.max(d.spend, d.income)), 1);

  return (
    <div className="cashflow-calendar-card" onClick={() => navigate('/calendar')}>
      <div className="card-header">
        <h2>This Week</h2>
        <button className="view-month-link">View Month</button>
      </div>

      <div className="calendar-week">
        {weekDays.map((day, index) => (
          <div key={index} className={`calendar-day ${day.isToday ? 'today' : ''}`}>
            <span className="day-label">{day.day}</span>
            <div className="day-number">{day.date}</div>

            <div className="day-bars">
              {day.spend > 0 && (
                <div
                  className="bar spend-bar"
                  style={{ height: `${(day.spend / maxAmount) * 40}px` }}
                  title={`Spend: $${day.spend}`}
                />
              )}
              {day.income > 0 && (
                <div
                  className="bar income-bar"
                  style={{ height: `${(day.income / maxAmount) * 40}px` }}
                  title={`Income: $${day.income}`}
                />
              )}
            </div>
          </div>
        ))}
      </div>

      <div className="calendar-legend">
        <div className="legend-item">
          <div className="legend-color spend"></div>
          <span>Spend</span>
        </div>
        <div className="legend-item">
          <div className="legend-color income"></div>
          <span>Income</span>
        </div>
      </div>
    </div>
  );
};

export default CashflowCalendar;