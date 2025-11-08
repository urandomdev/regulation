import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { financialService } from '../services/financialService';
import './CalendarWidget.css';

const CalendarWidget = () => {
  const navigate = useNavigate();
  const [currentDate, setCurrentDate] = useState(new Date());
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch transactions for the current month
  useEffect(() => {
    const fetchTransactions = async () => {
      setLoading(true);
      setError(null);
      try {
        const year = currentDate.getFullYear();
        const month = currentDate.getMonth();

        // Get first and last day of the month
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);

        // Format dates as YYYY-MM-DD
        const startDate = `${year}-${String(month + 1).padStart(2, '0')}-01`;
        const endDate = `${year}-${String(month + 1).padStart(2, '0')}-${String(lastDay.getDate()).padStart(2, '0')}`;

        // Fetch transactions for the month
        const result = await financialService.getTransactions({
          start: startDate,
          end: endDate,
          limit: 400, // Get all transactions for the month
        });

        // Transform transactions to match expected format
        const transformedTransactions = result.data.map(t => ({
          id: t.id,
          date: t.date,
          name: t.name || t.merchantName,
          category: t.category || 'Uncategorized',
          amount: t.isIncome ? -t.amount : t.amount, // Negative for income
          autoSaved: 0, // TODO: Get auto-saved amount from rules engine
        }));

        setTransactions(transformedTransactions);
      } catch (err) {
        console.error('Failed to fetch transactions:', err);
        setError('Failed to load transactions. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchTransactions();
  }, [currentDate]);

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
      const dayTransactions = transactions.filter(t => t.date === dateStr);

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
        {/* Error Message */}
        {error && (
          <div style={{
            padding: '1rem',
            marginBottom: '1rem',
            backgroundColor: '#fee',
            color: '#c33',
            borderRadius: '8px',
            border: '1px solid #fcc',
          }}>
            {error}
          </div>
        )}

        {/* Loading State */}
        {loading ? (
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '400px',
            gap: '1rem',
          }}>
            <div style={{
              width: '40px',
              height: '40px',
              border: '4px solid #e0e0e0',
              borderTop: '4px solid #4a90e2',
              borderRadius: '50%',
              animation: 'spin 1s linear infinite',
            }} />
            <p style={{ color: '#666' }}>Loading transactions...</p>
            <style>
              {`
                @keyframes spin {
                  0% { transform: rotate(0deg); }
                  100% { transform: rotate(360deg); }
                }
              `}
            </style>
          </div>
        ) : (
          <>
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
                ${transactions
                  .filter(t => t.amount < 0)
                  .reduce((sum, t) => sum + Math.abs(t.amount), 0)
                  .toFixed(2)}
              </div>
            </div>
            <div className="summary-item">
              <div className="summary-label">Total Spent</div>
              <div className="summary-value spend-value">
                ${transactions
                  .filter(t => t.amount > 0)
                  .reduce((sum, t) => sum + t.amount, 0)
                  .toFixed(2)}
              </div>
            </div>
            <div className="summary-item">
              <div className="summary-label">Auto-Saved</div>
              <div className="summary-value saved-value">
                ${transactions
                  .reduce((sum, t) => sum + t.autoSaved, 0)
                  .toFixed(2)}
              </div>
            </div>
          </div>
        </div>
          </>
        )}
      </div>
    </div>
  );
};

export default CalendarWidget;
