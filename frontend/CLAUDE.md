# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
```bash
npm run dev          # Start Vite development server on http://localhost:5173
npm run build        # Build for production (outputs to dist/)
npm run preview      # Preview production build locally
npm run lint         # Run ESLint to check for code issues
```

### Backend Dependency
This frontend requires a backend server running on `http://localhost:8080`. The Vite dev server proxies the following API routes to the backend:
- `/account` - Authentication and user management
- `/financial` - Financial data (accounts, transactions, cashflow)
- `/plaid` - Plaid integration
- `/advisor` - Financial advisor features
- `/notification` - Notifications

## Architecture Overview

### API Client Pattern
The app uses a custom SDK (`@urandomdev/regulation`) which is a local package dependency located at `../api/sdk/js`. The API client is initialized in `src/api/client.js` and provides two main API instances:
- `financialApi` - Financial operations (accounts, transactions, cashflow)
- `accountApi` - Authentication and user operations

API responses follow a tuple pattern: `[result, error]` where one is always null.

### Authentication Flow
- **Context Provider**: `src/contexts/AuthContext.jsx` wraps the entire app with authentication state
- **Private Routes**: `src/components/PrivateRoute.jsx` protects routes requiring authentication
- **Auth Methods**: `login()`, `signup()`, `logout()`, `checkAuth()`
- **Session Check**: On mount, the app calls `accountApi.me()` to verify existing session
- Unauthenticated users are redirected to `/login`

### Data Architecture
The app has a **dual data layer** during transition:

1. **Legacy Mock Data** (`src/data/store.js`):
   - Client-side state for rules, goals, and sample transactions
   - Used by older components not yet migrated to real API
   - Provides helper functions: `addRule()`, `updateRule()`, `deleteRule()`, `toggleRule()`, `addGoal()`, etc.

2. **Real API Service** (`src/services/financialService.js`):
   - Singleton service class wrapping the Financial API
   - Methods: `getAccounts()`, `getTransactions()`, `getCashflow()`, `getCurrentMonthCashflow()`
   - Includes caching layer and utility methods for balance/income/spending calculations
   - Converts cents (API) to dollars (UI) using `utils.centsToDollars()`

### Component Structure

**Pages** (Full-screen views accessible via routing):
- `Dashboard.jsx` - Main dashboard with overview cards
- `CalendarWidget.jsx` - Calendar view of financial events
- `GoalsWidget.jsx` - Full goals management page
- `CashflowWidget.jsx` - Detailed cashflow analysis
- `SavesHistoryWidget.jsx` - Transaction history
- `RulesWidget.jsx` - Automation rules management
- `Login.jsx`, `Signup.jsx` - Authentication pages

**Components** (Reusable UI elements):
- Dashboard cards: `FamilyManagement`, `BillingOverview`, `UpcomingMissions`, `CashflowSummary`
- Preview components: `GoalsPreview`, `RulesPreview`, `RecentSaves`, `TransactionHistory`
- `CashflowCalendar` - Calendar widget
- `PrivateRoute` - Route protection wrapper with loading state

### Styling Conventions
- Component-specific CSS files alongside each component (e.g., `Dashboard.css`, `BillingOverview.css`)
- Global styles in `src/index.css`
- Common/shared styles in `src/styles/common.css`
- CSS variables for theming (e.g., `var(--primary-blue)`, `var(--background)`, `var(--text-secondary)`)
- No external UI libraries - all styling is custom CSS

### Routing Structure
```
/                    → Dashboard (requires auth)
/login               → Login page
/signup              → Signup page
/calendar            → Calendar widget (requires auth)
/goals               → Goals management (requires auth)
/cashflow-details    → Cashflow details (requires auth)
/saves-history       → Transaction history (requires auth)
/rules               → Rules management (requires auth)
```

## Key Patterns

### API Error Handling
```javascript
const [result, error] = await accountApi.me();
if (error) {
  throw new Error(error.message || 'Operation failed');
}
// Use result
```

### Date Formatting
The `utils` object in `src/api/client.js` provides:
- `formatDate(date)` - Returns 'YYYY-MM-DD' string
- `getCurrentMonthRange()` - Returns `{ start, end }` for current month

### Widget Pattern
Full-page widgets follow this pattern:
- Back button to navigate to dashboard
- Header with title
- Main content area with data display
- Optional modals for forms (Add/Edit)
- State management for CRUD operations

### Currency Display
- Backend stores values in cents (integers)
- Frontend displays dollars (floats)
- Use `utils.centsToDollars(cents)` for conversion
- Transaction amounts shown as absolute values with `isIncome` flag

## Important Notes

- The app currently uses mock data in `src/data/store.js` for rules and goals, but real API data for accounts, transactions, and cashflow
- Korean language strings are present in auth error messages (e.g., "이메일 또는 비밀번호가 올바르지 않습니다")
- ESLint config allows unused variables starting with uppercase or underscore: `varsIgnorePattern: '^[A-Z_]'`
- The SDK dependency is a file path reference, not an npm package: `"@urandomdev/regulation": "file:../api/sdk/js"`
