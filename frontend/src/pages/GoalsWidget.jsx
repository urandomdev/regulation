import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userGoals, addGoal, updateGoal, deleteGoal } from '../data/store';
import './GoalsWidget.css';

const GoalsWidget = () => {
  const navigate = useNavigate();
  const [goals, setGoals] = useState([...userGoals]);
  const [showForm, setShowForm] = useState(false);
  const [editingGoal, setEditingGoal] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    targetAmount: '',
    type: 'Personal'
  });

  // Render goals list
  const renderGoalsList = () => {
    setGoals([...userGoals]);
  };

  // Calculate progress percentage
  const getProgressPercentage = (current, target) => {
    return Math.min(Math.round((current / target) * 100), 100);
  };

  // Handle delete
  const handleDelete = (goalId) => {
    if (window.confirm('Are you sure you want to delete this goal?')) {
      deleteGoal(goalId);
      renderGoalsList();
    }
  };

  // Handle edit
  const handleEdit = (goal) => {
    setEditingGoal(goal);
    setFormData({
      name: goal.name,
      targetAmount: goal.targetAmount.toString(),
      type: goal.type
    });
    setShowForm(true);
  };

  // Handle add contribution
  const handleAddContribution = (goalId) => {
    const amount = prompt('Enter contribution amount:');
    if (amount && !isNaN(amount) && parseFloat(amount) > 0) {
      const goal = userGoals.find(g => g.id === goalId);
      if (goal) {
        const newAmount = Math.min(
          goal.currentAmount + parseFloat(amount),
          goal.targetAmount
        );
        updateGoal(goalId, { currentAmount: newAmount });
        renderGoalsList();
      }
    }
  };

  // Handle form submit
  const handleSubmit = (e) => {
    e.preventDefault();

    const goalData = {
      name: formData.name,
      targetAmount: parseFloat(formData.targetAmount),
      type: formData.type
    };

    if (editingGoal) {
      // Update existing goal
      updateGoal(editingGoal.id, goalData);
    } else {
      // Add new goal
      addGoal(goalData);
    }

    // Reset form
    setFormData({ name: '', targetAmount: '', type: 'Personal' });
    setEditingGoal(null);
    setShowForm(false);
    renderGoalsList();
  };

  // Handle form cancel
  const handleCancel = () => {
    setFormData({ name: '', targetAmount: '', type: 'Personal' });
    setEditingGoal(null);
    setShowForm(false);
  };

  return (
    <div className="widget-page" id="goals-widget-page">
      <div className="widget-header">
        <button className="back-button" onClick={() => navigate('/')}>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        <h1>Your Goals</h1>
      </div>

      <div className="widget-content">
        <button
          className="add-button"
          onClick={() => setShowForm(true)}
        >
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add New Goal
        </button>

        <div id="goals-list-container" className="goals-list">
          {goals.map(goal => {
            const percentage = getProgressPercentage(goal.currentAmount, goal.targetAmount);
            const isCompleted = percentage === 100;

            return (
              <div key={goal.id} className={`goal-card ${isCompleted ? 'completed' : ''}`}>
                <div className="goal-header">
                  <div className="goal-title">
                    <h3 className="goal-name">{goal.name}</h3>
                    {goal.type === 'Family' && (
                      <span className="family-badge">👨‍👩‍👧‍👦</span>
                    )}
                  </div>
                  {isCompleted && (
                    <span className="completed-badge">✓ Completed</span>
                  )}
                </div>

                <div className="goal-progress">
                  <div className="progress-text">
                    <span className="current-amount">
                      ${goal.currentAmount.toFixed(2)}
                    </span>
                    <span className="separator"> / </span>
                    <span className="target-amount">
                      ${goal.targetAmount.toFixed(2)}
                    </span>
                  </div>

                  <div className="progress-bar-container">
                    <div
                      className="progress-bar-fill"
                      style={{ width: `${percentage}%` }}
                    ></div>
                  </div>

                  <div className="progress-percentage">
                    {percentage}% complete
                  </div>
                </div>

                <div className="goal-buttons">
                  {!isCompleted && (
                    <button
                      className="contribute-button"
                      onClick={() => handleAddContribution(goal.id)}
                    >
                      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                      Add Contribution
                    </button>
                  )}
                  <button
                    className="edit-button"
                    onClick={() => handleEdit(goal)}
                  >
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    Edit
                  </button>
                  <button
                    className="delete-button"
                    onClick={() => handleDelete(goal.id)}
                  >
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </div>
            );
          })}
        </div>

        {/* Add/Edit Goal Form Modal */}
        {showForm && (
          <div className="modal-overlay" onClick={handleCancel}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
              <h2>{editingGoal ? 'Edit Goal' : 'Add New Goal'}</h2>

              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="goal-name">Goal Name</label>
                  <input
                    type="text"
                    id="goal-name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="e.g., Save $500 this month"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="goal-target">Target Amount ($)</label>
                  <input
                    type="number"
                    id="goal-target"
                    value={formData.targetAmount}
                    onChange={(e) => setFormData({ ...formData, targetAmount: e.target.value })}
                    placeholder="e.g., 500"
                    min="1"
                    step="0.01"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="goal-type">Goal Type</label>
                  <select
                    id="goal-type"
                    value={formData.type}
                    onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                  >
                    <option value="Personal">Personal</option>
                    <option value="Family">Family</option>
                  </select>
                </div>

                <div className="form-buttons">
                  <button type="button" className="cancel-button" onClick={handleCancel}>
                    Cancel
                  </button>
                  <button type="submit" className="save-button">
                    {editingGoal ? 'Update Goal' : 'Save Goal'}
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

export default GoalsWidget;
