## Overview
**Who Are You To Spend?** is a simple, rule-based personal finance app.  
Users create straightforward rules that automatically move money into their savings account when they spend.  
The goal is not to restrict spending but to make saving a natural, automated reaction to everyday transactions.  
The app also provides a clear monthly cashflow summary showing income, expenses, savings, and net balance (NET).

---

## Target Users
1. **Teen**
   - Parent sets and manages rules.
   - Parent receives notifications for child’s transactions.
   - Designed for allowance and spending awareness.

2. **Me (Individual)**
   - User sets and manages their own rules.
   - Personalized rule recommendations based on spending habits.

3. **Family**
   - Each member can set personal rules.
   - Shared visibility into family spending and saving.
   - Family-wide saving goals available.

---

## Core Features

### 1. Account Linking
- Users can connect multiple accounts (Checking, Savings, Family Savings).
- One account can be designated as the **Auto-Save Account**.
- When a primary (main) checking account is linked, the app analyzes the last month of spending data and recommends starter rules.
  - Example: If Dining is a high category, suggest “Save 1.5× of every Dining expense.”

---

### 2. Spend Tracking
- Automatically imports transactions from linked accounts.
- Auto-categorizes spending into: Dining, Groceries, Transport, Shopping, Subscriptions, Entertainment, Bills, and Misc.
- Option to ignore small transactions (e.g., under $3).

---

### 3. Rules
Rules are the foundation of the app.  
Each rule follows a simple structure:

**Condition → Action**

When a transaction matches the condition, the specified amount is automatically moved to the savings account.

#### Rule Settings

**Conditions**
- Category (Dining, Shopping, etc.)
- Time of day (e.g., after 10 PM)
- Monthly total threshold (e.g., after $200 spent in Shopping)
- Transaction amount range (e.g., only apply to purchases over $10)

**Actions**
- Transfer a multiple of the spent amount (e.g., 1×, 1.5×, 2×)
- Transfer a fixed amount (e.g., $5)
- Ignore (no action)

#### Examples
- Dining purchases → Save 2× the amount.  
- Any spending after 10 PM → Save the same amount.  
- After $300 spent on Shopping → Save 30% of each following purchase.  
- Ignore transactions under $3.

Users can add up to 10 active rules and toggle them on or off.

---

### 4. Auto Transfers & Notifications
- When a rule condition is met, the app automatically transfers the calculated amount to the savings account.
- Notifications are minimal and purely informational:
  - “$15 moved to Savings.”
  - “$30 moved to Family Savings.”

---

### 5. Cashflow Summary
- Monthly summary showing:
  - Total Income  
  - Total Spend  
  - Auto-Save total  
  - NET (Income – Spend)
- Displays previous month comparison and top spending categories.
- Example:
```

Income: $3,200
Spend: $2,850
Auto-Save: $350
NET: +$350

```

---

### 6. Goals
- Users can set personal or family saving goals.
- Simple progress indicator showing percentage toward target.
- Example: “Save $500 this month – 70% complete.”
- If a goal is off track, the system may recommend an additional rule.

---

### 7. Family & Teen Management
- **Teen Mode**
- Parents can view child’s transactions and savings activity.
- Parents create or edit rules for the child’s account.
- **Family Mode**
- Each member has their own rules.
- Family can share summary data and track a shared goal (e.g., “Family Dining under $400/month”).
- Shared family rules available, such as:  
  “When total family Dining exceeds $400, move 2× excess to Family Savings.”

---

### 8. Onboarding Flow
1. Link primary checking account.
2. System analyzes last 30 days of spending.
3. Show category breakdown (e.g., Dining 25%, Shopping 20%).
4. Suggest starter rules (e.g., “Save 1.5× Dining spend”).
5. User selects 1–2 rules and starts tracking immediately.

---

### 9. Monthly Report
- Auto-generated summary at the end of each month.
- Includes:
- Income, Spending, Savings total
- NET
- Most active rules
- Goal completion
- Example output:  
“NET +$320 | Top Rule: Dining 2× Save | Goal: 80% complete”

---

### 10. Notification Policy
- Short, neutral, and numeric.
- Examples:
- “$12 moved to Savings.”
- “Goal reached.”
- “Family goal exceeded. $40 moved to Family Savings.”

---

## Tone and Style
- Minimal, neutral, and direct.
- No emotional or judgmental language.
- Focus on clarity, numbers, and results.

---

## One-line Summary
**Who Are You To Spend?** helps users define simple, universal money rules that automatically move funds to savings and clearly show their monthly NET.
