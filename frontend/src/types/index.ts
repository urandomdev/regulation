// Account types
export interface Account {
  id: string
  name: string
  type: 'checking' | 'savings' | 'other'
  mask?: string
  current_balance: number  // in cents
  available_balance?: number
  is_active: boolean
}

// Transaction types
export interface Transaction {
  id: string
  account_id: string
  amount: number  // cents (positive = debit, negative = credit)
  date: string
  name: string
  merchant_name?: string
  category: string
  pending: boolean
  payment_channel?: string
}

// Cashflow types
export interface Cashflow {
  total_income: number  // cents
  total_spend: number   // cents
  net: number           // cents
  start: string
  end: string
}

// Rule types
export type RuleCategory = 'dining' | 'groceries' | 'transport' | 'shopping' |
                           'subscriptions' | 'entertainment' | 'bills' | 'misc'

export type ActionType = 'multiply' | 'fixed'

export interface Rule {
  id: string
  name: string
  category: RuleCategory
  min_amount_cents?: number
  max_amount_cents?: number
  action_type: ActionType
  action_value: number
  target_account_id: string
  priority: number
  is_active: boolean
}

// User types
export interface User {
  id: string
  email: string
  nickname: string
}

// Auth types
export interface LoginRequest {
  email: string
  password: string
}

export interface SignupRequest {
  email: string
  password: string
  nickname: string
}
