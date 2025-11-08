import React from 'react';
import { useNavigate } from 'react-router-dom';
import { getMonthlyIncome, getMonthlySpend, getMonthlyAutoSaved } from '../data/store';
import './CashflowSummary.css';

const CashflowSummary = () => {
  const navigate = useNavigate();

  const income = getMonthlyIncome();
  const spend = getMonthlySpend();
  const saved = getMonthlyAutoSaved();
  const net = income - spend;

  return (
    <div className="cashflow-summary-card" onClick={() => navigate('/cashflow-details')}>
      <div className="card-header">
        <h2>Monthly Summary</h2>
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
      </div>

      <div className="summary-grid">
        <div className="summary-item income">
          <span className="summary-label">Income</span>
          <span className="summary-value">${income.toFixed(0)}</span>
        </div>
        <div className="summary-item spend">
          <span className="summary-label">Spend</span>
          <span className="summary-value">${spend.toFixed(0)}</span>
        </div>
        <div className="summary-item saved">
          <span className="summary-label">Auto-Saved</span>
          <span className="summary-value">${saved.toFixed(0)}</span>
        </div>
      </div>

      <div className="net-summary">
        <span className="net-label">NET</span>
        <span className={`net-value ${net >= 0 ? 'positive' : 'negative'}`}>
          {net >= 0 ? '+' : ''}${net.toFixed(0)}
        </span>
      </div>
    </div>
  );
};

export default CashflowSummary;