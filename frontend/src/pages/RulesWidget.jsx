import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { userRules, addRule, updateRule, deleteRule, toggleRule } from '../data/store';
import './RulesWidget.css';

const RulesWidget = () => {
  const navigate = useNavigate();
  const [rules, setRules] = useState([...userRules]);
  const [showForm, setShowForm] = useState(false);
  const [editingRule, setEditingRule] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    condition: '',
    action: ''
  });

  // Render rules list
  const renderRulesList = () => {
    setRules([...userRules]);
  };

  // Handle toggle switch
  const handleToggle = (ruleId) => {
    toggleRule(ruleId);
    renderRulesList();
  };

  // Handle delete
  const handleDelete = (ruleId) => {
    if (window.confirm('Are you sure you want to delete this rule?')) {
      deleteRule(ruleId);
      renderRulesList();
    }
  };

  // Handle edit
  const handleEdit = (rule) => {
    setEditingRule(rule);
    setFormData({
      name: rule.name,
      condition: rule.condition,
      action: rule.action
    });
    setShowForm(true);
  };

  // Handle form submit
  const handleSubmit = (e) => {
    e.preventDefault();

    if (editingRule) {
      // Update existing rule
      updateRule(editingRule.id, formData);
    } else {
      // Add new rule
      addRule(formData);
    }

    // Reset form
    setFormData({ name: '', condition: '', action: '' });
    setEditingRule(null);
    setShowForm(false);
    renderRulesList();
  };

  // Handle form cancel
  const handleCancel = () => {
    setFormData({ name: '', condition: '', action: '' });
    setEditingRule(null);
    setShowForm(false);
  };

  return (
    <div className="widget-page" id="rules-widget-page">
      <div className="widget-header">
        <button className="back-button" onClick={() =>{
            navigate('/');
        }
        }>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        <h1>Manage Rules</h1>
      </div>

      <div className="widget-content">
        <button
          className="add-button"
          onClick={() => {
              setShowForm(true);
          }
          }
        >
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add New Rule
        </button>

        <div id="rules-list-container" className="rules-list">
          {rules.map(rule => (
            <div key={rule.id} className={`rule-card ${!rule.isActive ? 'inactive' : ''}`}>
              <div className="rule-main">
                <div className="rule-info">
                  <h3 className="rule-name">{rule.name}</h3>
                  <p className="rule-details">
                    {rule.condition} → {rule.action}
                  </p>
                </div>

                <div className="rule-actions">
                  <label className="toggle-switch">
                    <input
                      type="checkbox"
                      checked={rule.isActive}
                      onChange={() => handleToggle(rule.id)}
                    />
                    <span className="toggle-slider"></span>
                  </label>
                </div>
              </div>

              <div className="rule-buttons">
                <button
                  className="edit-button"
                  onClick={() => handleEdit(rule)}
                >
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Edit
                </button>
                <button
                  className="delete-button"
                  onClick={() => handleDelete(rule.id)}
                >
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* Add/Edit Rule Form Modal */}
        {showForm && (
          <div className="modal-overlay" onClick={handleCancel}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
              <h2>{editingRule ? 'Edit Rule' : 'Add New Rule'}</h2>

              <form id="add-rule-form" onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="rule-name">Rule Name</label>
                  <input
                    type="text"
                    id="rule-name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="e.g., Dining 2x Save"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="rule-condition">Condition</label>
                  <input
                    type="text"
                    id="rule-condition"
                    value={formData.condition}
                    onChange={(e) => setFormData({ ...formData, condition: e.target.value })}
                    placeholder="e.g., Category: Dining"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="rule-action">Action</label>
                  <input
                    type="text"
                    id="rule-action"
                    value={formData.action}
                    onChange={(e) => setFormData({ ...formData, action: e.target.value })}
                    placeholder="e.g., Save 2x Amount"
                    required
                  />
                </div>

                <div className="form-buttons">
                  <button type="button" className="cancel-button" onClick={handleCancel}>
                    Cancel
                  </button>
                  <button type="submit" className="save-button">
                    {editingRule ? 'Update Rule' : 'Save Rule'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default RulesWidget;
