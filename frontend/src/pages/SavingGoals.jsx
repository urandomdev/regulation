import React from 'react';
import { useNavigate } from 'react-router-dom';
import './SavingGoals.css';

const SavingGoals = () => {
  const navigate = useNavigate();

  const goals = [
    {
      id: 1,
      icon: '💰',
      title: 'Emergency Fund',
      target: 5000,
      current: 3000,
      deadline: 'Dec 2025'
    },
    {
      id: 2,
      icon: '🏖️',
      title: 'Vacation Fund',
      target: 3000,
      current: 1050,
      deadline: 'Summer 2026'
    },
    {
      id: 3,
      icon: '🏠',
      title: 'Home Down Payment',
      target: 50000,
      current: 12000,
      deadline: 'Dec 2027'
    },
    {
      id: 4,
      icon: '🚗',
      title: 'New Car Fund',
      target: 15000,
      current: 4500,
      deadline: 'June 2026'
    },
  ];

  return (
    <div className="saving-goals-page">
      <div className="page-header">
        <button className="back-button" onClick={() => navigate('/')}>
          ← Back to Dashboard
        </button>
        <h1>Saving Goals</h1>
      </div>

      <div className="goals-grid">
        {goals.map((goal) => {
          const progress = (goal.current / goal.target) * 100;
          return (
            <div key={goal.id} className="goal-card">
              <div className="goal-icon">{goal.icon}</div>
              <h2>{goal.title}</h2>
              <div className="goal-amounts">
                <span className="current-amount">${goal.current.toLocaleString()}</span>
                <span className="separator">/</span>
                <span className="target-amount">${goal.target.toLocaleString()}</span>
              </div>
              <div className="goal-progress-bar">
                <div
                  className="goal-progress-fill"
                  style={{ width: `${progress}%` }}
                ></div>
              </div>
              <div className="goal-footer">
                <span className="progress-percentage">{progress.toFixed(0)}%</span>
                <span className="goal-deadline">{goal.deadline}</span>
              </div>
            </div>
          );
        })}
      </div>

      <div className="add-goal-section">
        <button className="add-goal-button">+ Add New Saving Goal</button>
      </div>
    </div>
  );
};

export default SavingGoals;