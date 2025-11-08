// Simulated backend data store
// This acts as our "fake backend" for dynamic components

// Sample user accounts
export let userAccounts = [
  { id: 'acc1', name: 'Main Checking', balance: 4200.50, type: 'Checking' },
  { id: 'acc2', name: 'My Savings', balance: 1500.00, type: 'Savings', isAutoSaveTarget: true },
  { id: 'acc3', name: 'Family Savings', balance: 500.00, type: 'Family Savings', isAutoSaveTarget: true }
];

// Sample user rules
export let userRules = [
  { id: 'rule1', name: 'Dining 2x Save', condition: 'Category: Dining', action: 'Save 2x Amount', isActive: true },
  { id: 'rule2', name: 'Late Night Spending', condition: 'Time: After 10 PM', action: 'Save 1x Amount', isActive: true },
  { id: 'rule3', name: 'Groceries', condition: 'Category: Groceries', action: 'Save $5 Fixed', isActive: false }
];

// Sample user goals
export let userGoals = [
  { id: 'goal1', name: 'Save $500 this month', currentAmount: 350, targetAmount: 500, type: 'Personal' },
  { id: 'goal2', name: 'Family Vacation', currentAmount: 200, targetAmount: 1000, type: 'Family' }
];

// Sample recent transactions
export let recentTransactions = [
  { id: 't1', date: '2025-11-08', name: 'Starbucks', category: 'Dining', amount: 6.50, autoSaved: 13.00 },
  { id: 't2', date: '2025-11-07', name: 'Amazon', category: 'Shopping', amount: 45.00, autoSaved: 0 },
  { id: 't3', date: '2025-11-07', name: 'Whole Foods', category: 'Groceries', amount: 120.00, autoSaved: 5.00 },
  { id: 't4', date: '2025-11-06', name: 'Coffee Shop', category: 'Dining', amount: 4.50, autoSaved: 9.00 },
  { id: 't5', date: '2025-11-05', name: 'Gas Station', category: 'Transport', amount: 55.00, autoSaved: 0 },
  { id: 't6', date: '2025-11-04', name: 'Restaurant', category: 'Dining', amount: 32.00, autoSaved: 64.00 },
  { id: 't7', date: '2025-11-03', name: 'Grocery Store', category: 'Groceries', amount: 85.00, autoSaved: 5.00 },
  { id: 't8', date: '2025-11-02', name: 'Online Shopping', category: 'Shopping', amount: 120.00, autoSaved: 0 },
  { id: 't9', date: '2025-11-01', name: 'Salary', category: 'Income', amount: -2500.00, autoSaved: 0 }
];

// Helper functions to manipulate the data store

// Rules
export const addRule = (rule) => {
  const newRule = {
    id: `rule${Date.now()}`,
    ...rule,
    isActive: true
  };
  userRules.push(newRule);
  return newRule;
};

export const updateRule = (id, updates) => {
  const index = userRules.findIndex(rule => rule.id === id);
  if (index !== -1) {
    userRules[index] = { ...userRules[index], ...updates };
    return userRules[index];
  }
  return null;
};

export const deleteRule = (id) => {
  userRules = userRules.filter(rule => rule.id !== id);
  return userRules;
};

export const toggleRule = (id) => {
  const rule = userRules.find(rule => rule.id === id);
  if (rule) {
    rule.isActive = !rule.isActive;
    return rule;
  }
  return null;
};

// Goals
export const addGoal = (goal) => {
  const newGoal = {
    id: `goal${Date.now()}`,
    currentAmount: 0,
    ...goal
  };
  userGoals.push(newGoal);
  return newGoal;
};

export const updateGoal = (id, updates) => {
  const index = userGoals.findIndex(goal => goal.id === id);
  if (index !== -1) {
    userGoals[index] = { ...userGoals[index], ...updates };
    return userGoals[index];
  }
  return null;
};

export const deleteGoal = (id) => {
  userGoals = userGoals.filter(goal => goal.id !== id);
  return userGoals;
};

// Transactions
export const addTransaction = (transaction) => {
  const newTransaction = {
    id: `t${Date.now()}`,
    ...transaction
  };
  recentTransactions.unshift(newTransaction);
  return newTransaction;
};

// Computed values
export const getTotalAutoSavedToday = () => {
  const today = new Date().toISOString().split('T')[0];
  return recentTransactions
    .filter(t => t.date === today)
    .reduce((sum, t) => sum + t.autoSaved, 0);
};

export const getActiveRulesCount = () => {
  return userRules.filter(rule => rule.isActive).length;
};

export const getTotalRulesCount = () => {
  return userRules.length;
};

export const getMonthlyIncome = () => {
  const thisMonth = new Date().toISOString().slice(0, 7);
  return Math.abs(recentTransactions
    .filter(t => t.date.startsWith(thisMonth) && t.amount < 0)
    .reduce((sum, t) => sum + t.amount, 0));
};

export const getMonthlySpend = () => {
  const thisMonth = new Date().toISOString().slice(0, 7);
  return recentTransactions
    .filter(t => t.date.startsWith(thisMonth) && t.amount > 0)
    .reduce((sum, t) => sum + t.amount, 0);
};

export const getMonthlyAutoSaved = () => {
  const thisMonth = new Date().toISOString().slice(0, 7);
  return recentTransactions
    .filter(t => t.date.startsWith(thisMonth))
    .reduce((sum, t) => sum + t.autoSaved, 0);
};
