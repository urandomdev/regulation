import React from 'react';
import { useNavigate } from 'react-router-dom';
import { userGoals } from '../data/store';
import './GoalsPreview.css';

const GoalsPreview = () => {
  const navigate = useNavigate();

  const goals = userGoals.slice(0, 2); // Show first 2 goals

  return (
    <div className="goals-preview-card" onClick={() => navigate('/goals')}>
      <div className="card-header">
        <h2>Goals</h2>
        <button className="view-all-link">View All</button>
      </div>

      <div className="goals-list">
        {goals.map((goal) => {
          const progress = (goal.currentAmount / goal.targetAmount) * 100;
          return (
            <div key={goal.id} className="goal-item">
              <div className="goal-info">
                <div className="goal-name">
                  {goal.type === 'Family' && <span className="family-badge">👨‍👩‍👧‍👦</span>}
                  {goal.name}
                </div>
                <div className="goal-amounts">
                  <span className="current">${goal.currentAmount.toFixed(0)}</span>
                  <span className="separator">/</span>
                  <span className="target">${goal.targetAmount.toFixed(0)}</span>
                </div>
              </div>
              <div className="goal-progress">
                <div className="progress-bar">
                  <div className="progress-fill" style={{ width: `${progress}%` }} />
                </div>
                <span className="progress-text">{Math.round(progress)}%</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default GoalsPreview;