export const CATEGORIES = [
  'dining',
  'groceries',
  'transport',
  'shopping',
  'subscriptions',
  'entertainment',
  'bills',
  'misc',
] as const

export const ACTION_TYPES = ['multiply', 'fixed'] as const

export const ACCOUNT_TYPES = ['checking', 'savings', 'other'] as const

export const CATEGORY_LABELS: Record<string, string> = {
  dining: 'Dining',
  groceries: 'Groceries',
  transport: 'Transport',
  shopping: 'Shopping',
  subscriptions: 'Subscriptions',
  entertainment: 'Entertainment',
  bills: 'Bills',
  misc: 'Miscellaneous',
}
