import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userRules, getActiveRulesCount, getTotalRulesCount } from '../data/store';
import './RulesPreview.css';

const RulesPreview = () => {
    const navigate = useNavigate();
    const [rules, setRules] = useState(userRules.slice(0, 3)); 

    const activeRules = getActiveRulesCount();
    const totalRules = getTotalRulesCount();

    const toggleRule = (id) => {
        setRules((prev) =>
            prev.map((rule) =>
                rule.id === id ? { ...rule, isActive: !rule.isActive } : rule
            )
        );
    };

    return (
        <div className="rules-preview-card">
            <div className="card-header">
                <h2>Active Rules</h2>
                <span className="rules-count">
          {activeRules}/{totalRules}
        </span>
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

                        <div
                            className={`rule-toggle ${rule.isActive ? 'active' : ''}`}
                            onClick={() => toggleRule(rule.id)}
                        >
                            <div className="toggle-indicator"></div>
                        </div>
                    </div>
                ))}
            </div>

            <button
                className="add-rule-button"
                onClick={(e) => {
                    e.stopPropagation();
                    navigate('/rules');
                }}
            >
                + Add New Rule
            </button>
        </div>
    );
};

export default RulesPreview;
