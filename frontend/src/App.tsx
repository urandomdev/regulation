import { SignupForm } from './components/SignupForm'
import { ApiTester } from './components/ApiTester'
import './App.css'

function App() {
  return (
    <div style={{ padding: '20px' }}>
      <h1 style={{ textAlign: 'center' }}>Spending Tracker - API Integration Demo</h1>

      <div style={{
        display: 'grid',
        gridTemplateColumns: '1fr 1fr',
        gap: '40px',
        maxWidth: '1200px',
        margin: '0 auto'
      }}>
        <div>
          <SignupForm />
        </div>

        <div>
          <ApiTester />
        </div>
      </div>
    </div>
  )
}

export default App
