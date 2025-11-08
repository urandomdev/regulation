import React from 'react';
import './TransactionHistory.css';

const TransactionHistory = ({ isVisible }) => {
  const transactions = [
    { id: 1, date: '2025-11-01', description: 'Monthly Subscription', amount: '-$29.99' },
    { id: 2, date: '2025-11-05', description: 'Payment Received', amount: '+$150.00' },
    { id: 3, date: '2025-11-07', description: 'Service Upgrade', amount: '-$49.99' },
  ];

  if (!isVisible) return null;

  return (
    <div className="transaction-history-container">
      <h3>Transaction History</h3>
      <div className="transactions-list">
        {transactions.map((transaction) => (
          <div key={transaction.id} className="transaction-item">
            <div className="transaction-info">
              <span className="transaction-date">{transaction.date}</span>
              <span className="transaction-description">{transaction.description}</span>
            </div>
            <span className={`transaction-amount ${transaction.amount.startsWith('+') ? 'positive' : 'negative'}`}>
              {transaction.amount}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TransactionHistory;
