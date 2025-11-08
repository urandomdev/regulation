import { Routes, Route } from 'react-router'
import Home from './pages/Home'
import Login from './pages/Login'
import Signup from './pages/Signup'
import Dashboard from './pages/Dashboard'
import Accounts from './pages/Accounts'
import AccountDetail from './pages/AccountDetail'
import Rules from './pages/Rules'
import NotFound from './pages/NotFound'
import Onboarding from './pages/Onboarding'
import DisconnectAccounts from './pages/DisconnectAccounts'
import Recommendations from './pages/Recommendations'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/accounts" element={<Accounts />} />
      <Route path="/accounts/:accountId" element={<AccountDetail />} />
      <Route path="/accounts/disconnect" element={<DisconnectAccounts />} />
      <Route path="/onboarding" element={<Onboarding />} />
      <Route path="/recommendations" element={<Recommendations />} />
      <Route path="/rules" element={<Rules />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  )
}

export default App
