import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import PrivateRoute from './components/PrivateRoute';
import Dashboard from './pages/Dashboard';
import CalendarWidget from './pages/CalendarWidget';
import GoalsWidget from './pages/GoalsWidget';
import CashflowWidget from './pages/CashflowWidget';
import SavesHistoryWidget from './pages/SavesHistoryWidget';
import RulesWidget from './pages/RulesWidget';
import Login from './pages/Login';
import Signup from './pages/Signup';
import './App.css';

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="app">
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
            <Route path="/" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
            <Route path="/calendar" element={<PrivateRoute><CalendarWidget /></PrivateRoute>} />
            <Route path="/goals" element={<PrivateRoute><GoalsWidget /></PrivateRoute>} />
            <Route path="/cashflow-details" element={<PrivateRoute><CashflowWidget /></PrivateRoute>} />
            <Route path="/saves-history" element={<PrivateRoute><SavesHistoryWidget /></PrivateRoute>} />
            <Route path="/rules" element={<PrivateRoute><RulesWidget /></PrivateRoute>} />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;