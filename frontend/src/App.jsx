import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import CalendarWidget from './pages/CalendarWidget';
import GoalsWidget from './pages/GoalsWidget';
import CashflowWidget from './pages/CashflowWidget';
import SavesHistoryWidget from './pages/SavesHistoryWidget';
import RulesWidget from './pages/RulesWidget';
import './App.css';

function App() {
  return (
    <Router>
      <div className="app">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/calendar" element={<CalendarWidget />} />
          <Route path="/goals" element={<GoalsWidget />} />
          <Route path="/cashflow-details" element={<CashflowWidget />} />
          <Route path="/saves-history" element={<SavesHistoryWidget />} />
          <Route path="/rules" element={<RulesWidget />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;