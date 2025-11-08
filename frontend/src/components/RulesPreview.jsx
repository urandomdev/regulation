import React from 'react';
import { useNavigate } from 'react-router-dom';
import { userRules, getActiveRulesCount, getTotalRulesCount } from '../data/store';
import './RulesPreview.css';

const RulesPreview = () => {
  const navigate = useNavigate();

  const rules = userRules.slice(0, 3); // Show first 3 rules
  const activeRules = getActiveRulesCount();
  const totalRules = getTotalRulesCount();

  return (
    <div className="rules-preview-card" onClick={() => navigate('/rules')}>
      <div className="card-header">
        <h2>Active Rules</h2>
        <span className="rules-count">{activeRules}/{totalRules}</span>
      </div>

      <div className="rules-list">
        {rules.map((rule) => (
          <div key={rule.id} className="rule-item">
            <div className="rule-info">
              <div className="rule-name">{rule.name}</div>
              <div className="rule-stats">
                {rule.condition} → {rule.action}
              </div>
            </div>
            <div className={`rule-toggle ${rule.isActive ? 'active' : ''}`}>
              <div className="toggle-indicator"></div>
            </div>
          </div>
        ))}
      </div>

      <button className="add-rule-button" onClick={(e) => {
        e.stopPropagation();
        navigate('/rules');
      }}>
        + Add New Rule
      </button>
    </div>
  );
};

export default RulesPreview;