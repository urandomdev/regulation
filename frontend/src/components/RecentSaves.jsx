import React from 'react';
import { useNavigate } from 'react-router-dom';
import { recentTransactions, getTotalAutoSavedToday } from '../data/store';
import './RecentSaves.css';

const RecentSaves = () => {
  const navigate = useNavigate();

  // Get recent saves (transactions with autoSaved > 0)
  const recentSaves = recentTransactions
    .filter(t => t.autoSaved > 0)
    .slice(0, 3)
    .map(t => {
      const txDate = new Date(t.date);
      const now = new Date();
      const diffInHours = Math.floor((now - txDate) / (1000 * 60 * 60));

      let timeAgo;
      if (diffInHours < 1) {
        timeAgo = 'Just now';
      } else if (diffInHours < 24) {
        timeAgo = `${diffInHours}h ago`;
      } else {
        const diffInDays = Math.floor(diffInHours / 24);
        timeAgo = `${diffInDays}d ago`;
      }

      return {
        id: t.id,
        description: t.name,
        amount: t.autoSaved,
        account: 'My Savings',
        time: timeAgo
      };
    });

  const totalSavedToday = getTotalAutoSavedToday();

  return (
    <div className="recent-saves-card" onClick={() => navigate('/saves-history')}>
      <div className="card-header">
        <h2>Recent Auto-Saves</h2>
        <button className="view-all-link">View All</button>
      </div>

      <div className="saves-list">
        {recentSaves.map((save) => (
          <div key={save.id} className="save-item">
            <div className="save-info">
              <div className="save-description">{save.description}</div>
              <div className="save-meta">
                <span className="save-account">{save.account}</span>
                <span className="save-time">{save.time}</span>
              </div>
            </div>
            <div className="save-amount">+${save.amount}</div>
          </div>
        ))}
      </div>

      <div className="total-saved">
        <span className="total-label">Total Auto-Saved Today</span>
        <span className="total-amount">${totalSavedToday.toFixed(2)}</span>
      </div>
    </div>
  );
};

export default RecentSaves;