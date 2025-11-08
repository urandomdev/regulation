import React from 'react';
import { useNavigate } from 'react-router-dom';
import CashflowCalendar from '../components/CashflowCalendar';
import GoalsPreview from '../components/GoalsPreview';
import RecentSaves from '../components/RecentSaves';
import CashflowSummary from '../components/CashflowSummary';
import RulesPreview from '../components/RulesPreview';
import './Dashboard.css';

const Dashboard = () => {
  const navigate = useNavigate();

  return (
    <div className="dashboard-container">
      {/* Header */}
      <div className="dashboard-header">
        <h1>Hello, Alex!</h1>
        <button className="header-add-button" onClick={() => navigate('/add-rule')}>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      {/* Content */}
      <div className="dashboard-content">
        <CashflowSummary />
        <CashflowCalendar />
        <GoalsPreview />
        <RecentSaves />
        <RulesPreview />
      </div>

      {/* Floating Action Button */}
      <button className="fab" onClick={() => navigate('/add-rule')}>
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </div>
  );
};

export default Dashboard;