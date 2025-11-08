import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { recentTransactions } from '../data/store';
import './CalendarWidget.css';

const CalendarWidget = () => {
  const navigate = useNavigate();
  const [currentDate, setCurrentDate] = useState(new Date());

  // Get calendar data for current month
  const getCalendarData = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();

    // Get first day of month and number of days
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay(); // 0 = Sunday

    // Create array of day objects
    const days = [];

    // Add empty cells for days before month starts
    for (let i = 0; i < startingDayOfWeek; i++) {
      days.push({ isEmpty: true });
    }

    // Add actual days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      const dayTransactions = recentTransactions.filter(t => t.date === dateStr);

      const hasSpending = dayTransactions.some(t => t.amount > 0);
      const hasIncome = dayTransactions.some(t => t.amount < 0);
      const totalSpend = dayTransactions
        .filter(t => t.amount > 0)
        .reduce((sum, t) => sum + t.amount, 0);
      const totalIncome = Math.abs(dayTransactions
        .filter(t => t.amount < 0)
        .reduce((sum, t) => sum + t.amount, 0));
      const totalAutoSaved = dayTransactions.reduce((sum, t) => sum + t.autoSaved, 0);

      const isToday =
        day === new Date().getDate() &&
        month === new Date().getMonth() &&
        year === new Date().getFullYear();

      days.push({
        day,
        dateStr,
        hasSpending,
        hasIncome,
        totalSpend,
        totalIncome,
        totalAutoSaved,
        transactions: dayTransactions,
        isToday
      });
    }

    return days;
  };

  // Navigate to previous month
  const previousMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1));
  };

  // Navigate to next month
  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1));
  };

  // Format month and year
  const monthYear = currentDate.toLocaleDateString('en-US', {
    month: 'long',
    year: 'numeric'
  });

  const calendarDays = getCalendarData();

  return (
    <div className="widget-page" id="calendar-widget-page">
      <div className="widget-header">
        <button className="back-button" onClick={() => navigate('/')}>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        <h1>Monthly Cashflow</h1>
      </div>

      <div className="widget-content">
        {/* Month Navigation */}
        <div className="calendar-header">
          <button className="month-nav-button" onClick={previousMonth}>
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h2 className="month-title">{monthYear}</h2>
          <button className="month-nav-button" onClick={nextMonth}>
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        {/* Calendar Grid */}
        <div className="calendar-container">
          {/* Day names */}
          <div className="calendar-grid">
            {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(dayName => (
              <div key={dayName} className="day-name">
                {dayName}
              </div>
            ))}

            {/* Calendar days */}
            {calendarDays.map((dayData, index) => {
              if (dayData.isEmpty) {
                return <div key={`empty-${index}`} className="calendar-day empty"></div>;
              }

              return (
                <div
                  key={dayData.dateStr}
                  className={`calendar-day ${dayData.isToday ? 'today' : ''} ${
                    dayData.transactions.length > 0 ? 'has-transactions' : ''
                  }`}
                >
                  <div className="day-number">{dayData.day}</div>

                  {/* Transaction indicators */}
                  {(dayData.hasSpending || dayData.hasIncome) && (
                    <div className="transaction-indicators">
                      {dayData.hasSpending && (
                        <div className="indicator spending-dot" title={`Spent: $${dayData.totalSpend.toFixed(2)}`}></div>
                      )}
                      {dayData.hasIncome && (
                        <div className="indicator income-dot" title={`Income: $${dayData.totalIncome.toFixed(2)}`}></div>
                      )}
                    </div>
                  )}

                  {/* Show amounts for days with transactions */}
                  {dayData.transactions.length > 0 && (
                    <div className="day-amounts">
                      {dayData.totalSpend > 0 && (
                        <div className="amount spend-amount">
                          ${dayData.totalSpend.toFixed(0)}
                        </div>
                      )}
                      {dayData.totalAutoSaved > 0 && (
                        <div className="amount saved-amount">
                          +${dayData.totalAutoSaved.toFixed(0)}
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* Legend */}
        <div className="calendar-legend">
          <h3>Legend</h3>
          <div className="legend-items">
            <div className="legend-item">
              <div className="legend-dot spending-dot"></div>
              <span>Spending</span>
            </div>
            <div className="legend-item">
              <div className="legend-dot income-dot"></div>
              <span>Income</span>
            </div>
            <div className="legend-item">
              <div className="legend-indicator today-indicator"></div>
              <span>Today</span>
            </div>
          </div>
        </div>

        {/* Month Summary */}
        <div className="month-summary">
          <h3>This Month Summary</h3>
          <div className="summary-grid">
            <div className="summary-item">
              <div className="summary-label">Total Income</div>
              <div className="summary-value income-value">
                ${recentTransactions
                  .filter(t => {
                    const tDate = new Date(t.date);
                    return (
                      t.amount < 0 &&
                      tDate.getMonth() === currentDate.getMonth() &&
                      tDate.getFullYear() === currentDate.getFullYear()
                    );
                  })
                  .reduce((sum, t) => sum + Math.abs(t.amount), 0)
                  .toFixed(2)}
              </div>
            </div>
            <div className="summary-item">
              <div className="summary-label">Total Spent</div>
              <div className="summary-value spend-value">
                ${recentTransactions
                  .filter(t => {
                    const tDate = new Date(t.date);
                    return (
                      t.amount > 0 &&
                      tDate.getMonth() === currentDate.getMonth() &&
                      tDate.getFullYear() === currentDate.getFullYear()
                    );
                  })
                  .reduce((sum, t) => sum + t.amount, 0)
                  .toFixed(2)}
              </div>
            </div>
            <div className="summary-item">
              <div className="summary-label">Auto-Saved</div>
              <div className="summary-value saved-value">
                ${recentTransactions
                  .filter(t => {
                    const tDate = new Date(t.date);
                    return (
                      tDate.getMonth() === currentDate.getMonth() &&
                      tDate.getFullYear() === currentDate.getFullYear()
                    );
                  })
                  .reduce((sum, t) => sum + t.autoSaved, 0)
                  .toFixed(2)}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CalendarWidget;