import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { financialService } from '../services/financialService';
import './CashflowSummary.css';

const CashflowSummary = () => {
  const navigate = useNavigate();
  const [cashflow, setCashflow] = useState({
    income: 0,
    spend: 0,
    net: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCashflow = async () => {
      try {
        setLoading(true);
        const data = await financialService.getCurrentMonthCashflow();
        setCashflow(data);
      } catch (err) {
        console.error('Failed to fetch cashflow:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchCashflow();
  }, []);

  if (loading) {
    return (
      <div className="cashflow-summary-card">
        <div className="card-header">
          <h2>Monthly Summary</h2>
        </div>
        <div className="loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="cashflow-summary-card">
        <div className="card-header">
          <h2>Monthly Summary</h2>
        </div>
        <div className="error">Error: {error}</div>
      </div>
    );
  }

  const { income, spend, net } = cashflow;
  // TODO: Auto-saved amount - need backend API support
  const saved = 0;

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
          <span className="summary-value">${Math.abs(income).toFixed(0)}</span>
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