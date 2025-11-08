import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { financialService } from '../services/financialService';
import './RecentSaves.css';

const RecentSaves = () => {
  const navigate = useNavigate();
  const [recentSaves, setRecentSaves] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchRecentTransactions = async () => {
      try {
        setLoading(true);
        const transactions = await financialService.getRecentTransactions(10);

        // TODO: Backend needs to support autoSaved field
        // For now, showing recent transactions as placeholder
        const saves = transactions.slice(0, 3).map(t => {
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
            amount: 0, // TODO: Use t.autoSaved when backend supports it
            account: 'My Savings',
            time: timeAgo
          };
        });

        setRecentSaves(saves);
      } catch (err) {
        console.error('Failed to fetch recent transactions:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchRecentTransactions();
  }, []);

  if (loading) {
    return (
      <div className="recent-saves-card">
        <div className="card-header">
          <h2>Recent Auto-Saves</h2>
        </div>
        <div className="loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="recent-saves-card">
        <div className="card-header">
          <h2>Recent Auto-Saves</h2>
        </div>
        <div className="error">Error: {error}</div>
      </div>
    );
  }

  const totalSavedToday = 0; // TODO: Calculate from backend when autoSaved is supported

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