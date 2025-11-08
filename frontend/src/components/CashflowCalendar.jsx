import React from 'react';
import { useNavigate } from 'react-router-dom';
import './CashflowCalendar.css';

const CashflowCalendar = () => {
  const navigate = useNavigate();

  const weekDays = [
    { day: 'MON', date: 21, spend: 45, income: 0 },
    { day: 'TUE', date: 22, spend: 32, income: 0 },
    { day: 'WED', date: 23, spend: 28, income: 0, isToday: true },
    { day: 'THU', date: 24, spend: 51, income: 0 },
    { day: 'FRI', date: 25, spend: 38, income: 0 },
    { day: 'SAT', date: 26, spend: 72, income: 0 },
    { day: 'SUN', date: 27, spend: 19, income: 150 },
  ];

  const maxAmount = Math.max(...weekDays.map(d => Math.max(d.spend, d.income)));

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