import React from 'react';
import { useNavigate } from 'react-router-dom';
import './WidgetPage.css';

const CashflowWidget = () => {
  const navigate = useNavigate();

  return (
    <div className="widget-page">
      <div className="widget-header">
        <button className="back-button" onClick={() => navigate('/')}>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        <h1>Cashflow Details</h1>
      </div>

      <div className="widget-content">
        <p className="placeholder-text">Detailed cashflow breakdown and analytics will be displayed here.</p>
      </div>
    </div>
  );
};

export default CashflowWidget;