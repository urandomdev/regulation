import React, { useState } from 'react';
import TransactionHistory from './TransactionHistory';
import './BillingOverview.css';

const BillingOverview = () => {
  const [showTransactionHistory, setShowTransactionHistory] = useState(false);

  const handleClick = () => {
    setShowTransactionHistory(!showTransactionHistory);
  };

  return (
    <div className="billing-overview-container">
      <div className="billing-header" onClick={handleClick}>
        <h2>Billing Overview</h2>
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
      </div>
      <div className="billing-content">
        <span className="billing-label">Unpaid Invoices</span>
        <p className="billing-amount">$1,250.75</p>
      </div>
      <TransactionHistory isVisible={showTransactionHistory} />
    </div>
  );
};

export default BillingOverview;