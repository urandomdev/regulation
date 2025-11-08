# React UI/UX Application

A React application with Family Management, Billing Overview, and Saving Goals features.

## Features Implemented

### 1. Family Management Component
- **Renamed from**: "pedding controll" (as requested)
- **Location**: `src/components/FamilyManagement.jsx`
- Displays family management interface with clean design

### 2. Billing Overview with Transaction History
- **Location**: `src/components/BillingOverview.jsx`
- **Interaction**: Clicking on "Billing Overview" toggles the Transaction History view
- Shows billing summary with current balance, next payment, and payment method
- Transaction History displays as an expandable section below the billing overview

### 3. Upcoming Missions with Navigation
- **Location**: `src/components/UpcomingMissions.jsx`
- **Navigation**: Clicking on "Upcoming Missions" navigates to the Saving Goals page
- Displays progress on current saving missions
- Visual hover effects indicate clickability

### 4. Saving Goals Page
- **Location**: `src/pages/SavingGoals.jsx`
- Full page dedicated to all saving goals
- Displays goal progress, target amounts, and deadlines
- Back button to return to dashboard

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── FamilyManagement.jsx
│   │   ├── FamilyManagement.css
│   │   ├── BillingOverview.jsx
│   │   ├── BillingOverview.css
│   │   ├── TransactionHistory.jsx
│   │   ├── TransactionHistory.css
│   │   ├── UpcomingMissions.jsx
│   │   └── UpcomingMissions.css
│   ├── pages/
│   │   ├── Dashboard.jsx
│   │   ├── Dashboard.css
│   │   ├── SavingGoals.jsx
│   │   └── SavingGoals.css
│   ├── App.jsx
│   ├── App.css
│   └── main.jsx
├── package.json
└── vite.config.js
```

## Installation & Running

### Prerequisites
- Node.js (v14 or higher)
- npm or yarn

### Steps

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies (if not already installed):
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

4. Open your browser and navigate to:
```
http://localhost:5173/
```

## Usage

### Dashboard
The main dashboard displays three main components:
- **Family Management**: Manage family members and accounts
- **Billing Overview**: View billing summary (click to expand Transaction History)
- **Upcoming Missions**: View saving progress (click to navigate to full Saving Goals page)

### Navigation
- Click on **"Upcoming Missions"** to navigate to the Saving Goals page
- Click on **"Billing Overview"** header to toggle Transaction History
- Use the **"Back to Dashboard"** button on the Saving Goals page to return

## Design Philosophy

- **Clean & Minimal**: Simple, modern design with card-based layouts
- **Consistent Styling**: Uniform color scheme and spacing throughout
- **Responsive**: Adapts to different screen sizes
- **Interactive**: Hover effects and transitions for better UX
- **Accessible**: Clear typography and color contrast

## Technologies Used

- **React**: UI library
- **React Router**: Client-side routing
- **Vite**: Build tool and development server
- **CSS3**: Styling with modern features (Grid, Flexbox, Transitions)

## Build for Production

```bash
npm run build
```

The production-ready files will be in the `dist/` directory.

## Notes

- All components maintain the existing design system
- No external UI libraries were used to maintain full design control
- All interactions are implemented with React state management
- Routing is handled client-side with React Router
