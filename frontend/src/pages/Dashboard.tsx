import { useCallback, useEffect, useState } from 'react'
import { UUID } from '@deltalaboratory/uuid'
import { Plus, CreditCard, ChevronRight, LogOut } from 'lucide-react'
import { useNavigate } from 'react-router'
import { useAuth } from '@/hooks/useAuth'
import { api } from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { PlaidLinkButton } from '@/components/PlaidLink'
import { PageShell, PageHeader, PageBody } from '@/components/layout/Page'

interface User {
  id: UUID
  email: string
  nickname?: string
}

interface Account {
  id: UUID
  name: string
  type: string
  mask?: string | null
  current_balance: bigint
  available_balance?: bigint | null
  is_active: boolean
}

interface Rule {
  id: UUID
  name: string
  category: string
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  action_type: string
  action_value: number
  target_account_id: UUID
  priority: bigint
  is_active: boolean
  execution_count: bigint
  total_saved_cents: bigint
  created_at: Date
  updated_at: Date
}

interface Cashflow {
  total_income: bigint
  total_spend: bigint
  net: bigint
  start: Date
  end: Date
}

interface Transaction {
  id: UUID
  account_id: UUID
  amount: bigint
  date: string
  name: string
  merchant_name?: string
  category: string
  pending: boolean
}

const Dashboard = () => {
  const [user, setUser] = useState<User | null>(null)
  const [accounts, setAccounts] = useState<Account[]>([])
  const [rules, setRules] = useState<Rule[]>([])
  const [cashflow, setCashflow] = useState<Cashflow | null>(null)
  const [recentTransactions, setRecentTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { logout } = useAuth()
  const navigate = useNavigate()

  const fetchData = useCallback(async () => {
    // Fetch user
    const [userData, userErr] = await api.account.me()
    if (userErr) {
      setError(userErr.message || 'Failed to fetch user')
      setLoading(false)
      return
    }
    setUser(userData as User)

    // Fetch accounts
    const [accountsData, accountsErr] = await api.financial.getAccounts()
    if (!accountsErr && accountsData) {
      setAccounts(accountsData.accounts)

      if (accountsData.accounts.length === 0) {
        navigate('/onboarding', { replace: true })
        return
      }
    }

    // Fetch rules
    const [rulesData, rulesErr] = await api.rule.listRules()
    if (!rulesErr && rulesData) {
      setRules(rulesData.rules)
    }

    // Fetch cashflow for current month
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)

    const [cashflowData, cashflowErr] = await api.financial.getCashflow({
      start: startOfMonth,
      end: endOfMonth,
    })
    if (!cashflowErr && cashflowData) {
      setCashflow(cashflowData)
    }

    // Fetch recent transactions (last 5)
    const [transactionsData, transactionsErr] = await api.financial.getTransactions({
      limit: BigInt(5),
      offset: BigInt(0),
    })
    if (!transactionsErr && transactionsData) {
      // Convert SDK transactions to local type
      const convertedTransactions: Transaction[] = (transactionsData.transactions || []).map(tx => ({
        id: tx.id,
        account_id: tx.account_id,
        amount: tx.amount,
        date: tx.date.toISOString(),
        name: tx.name,
        merchant_name: tx.merchant_name ?? undefined,
        category: tx.category,
        pending: tx.pending
      }))
      setRecentTransactions(convertedTransactions)
    }

    setLoading(false)
  }, [navigate])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  if (loading) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <p className="text-neutral-500">Loading...</p>
        </PageBody>
      </PageShell>
    )
  }

  if (error) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <p className="text-destructive">Error: {error}</p>
        </PageBody>
      </PageShell>
    )
  }

  // Calculate cashflow stats
  const income = cashflow ? Number(cashflow.total_income) / 100 : 0
  const spend = cashflow ? Number(cashflow.total_spend) / 100 : 0
  const net = cashflow ? Number(cashflow.net) / 100 : 0

  // Calculate auto save from rules
  const autoSave = rules.reduce((sum, rule) => sum + Number(rule.total_saved_cents) / 100, 0)

  // Goal (placeholder for now, could be fetched from user settings)
  const goal = 500
  const goalProgress = goal > 0 ? Math.min(Math.round((autoSave / goal) * 100), 100) : 0

  return (
    <PageShell>
      <PageHeader
        title={user?.nickname || 'Dashboard'}
        action={(
          <Button
            variant="ghost"
            size="icon"
            className="text-neutral-400 hover:text-neutral-600 transition-colors"
            onClick={handleLogout}
          >
            <LogOut className="h-5 w-5" />
          </Button>
        )}
      />

      <PageBody className="space-y-6">
        {/* Monthly Summary */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5 space-y-3">
            <div className="flex justify-between items-baseline">
              <h2 className="text-lg font-semibold tracking-tight">This Month</h2>
              <span className="text-sm text-neutral-500">Goal ${goal}</span>
            </div>
            <div className="text-3xl font-bold tracking-tight text-black">NET +${net}</div>
            <div className="grid grid-cols-3 gap-3 text-sm">
              <KPI label="Income" value={`${income}`} />
              <KPI label="Spend" value={`${spend}`} />
              <KPI label="Auto Save" value={`${autoSave}`} />
            </div>
            <Progress value={goalProgress} className="h-1.5 mt-2 rounded-full bg-gray-100 [&>div]:bg-teal-500" />
          </CardContent>
        </Card>

        {/* Accounts Preview */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5 space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold text-neutral-700">Accounts</h3>
              <PlaidLinkButton
                variant="ghost"
                size="sm"
                className="text-sm text-neutral-500 hover:text-neutral-700 transition-colors"
                onSuccess={fetchData}
              >
                <Plus className="h-4 w-4 mr-1" />Add
              </PlaidLinkButton>
            </div>
            <div className="space-y-2">
              {accounts.length > 0 ? (
                accounts.map((a) => (
                  <button
                    key={a.id.toString()}
                    className="flex w-full items-center justify-between rounded-lg bg-white border border-gray-100 p-3 text-sm hover:border-gray-200 transition-colors"
                    onClick={() => navigate(`/accounts/${a.id.toString()}`)}
                  >
                    <div className="flex items-center gap-2">
                      <CreditCard className="h-4 w-4 text-gray-400" />
                      <span className="font-medium">{a.name}</span>
                      {a.mask && <span className="text-neutral-500">••{a.mask}</span>}
                    </div>
                    <div className="flex items-center gap-1 font-semibold text-neutral-700">
                      <span>${(Number(a.current_balance) / 100).toFixed(2)}</span>
                      <ChevronRight className="h-4 w-4 text-neutral-300" />
                    </div>
                  </button>
                ))
              ) : (
                <div className="text-center py-4 text-sm text-neutral-500">
                  No accounts connected yet
                </div>
              )}
            </div>
            <Button
              variant="outline"
              size="sm"
              className="w-full mt-2 rounded-lg border-gray-200 text-neutral-700 hover:border-gray-300 hover:text-neutral-900 transition-colors"
              onClick={() => navigate('/accounts')}
            >
              View All Accounts
            </Button>
          </CardContent>
        </Card>

        {/* Rules Preview */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5 space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold text-neutral-700">My Rules</h3>
              <Button
                size="sm"
                className="rounded-lg text-sm bg-gray-100 text-neutral-800 hover:bg-gray-200 transition-colors"
                onClick={() => navigate('/rules')}
              >
                <Plus className="h-4 w-4 mr-1" />Add
              </Button>
            </div>
            <div className="space-y-1">
              {rules.length > 0 ? (
                rules.slice(0, 3).map((rule) => (
                  <Rule key={rule.id.toString()} text={rule.name} isActive={rule.is_active} />
                ))
              ) : (
                <div className="text-center py-4 text-sm text-neutral-500">
                  No rules created yet
                </div>
              )}
            </div>
            <Button
              variant="outline"
              size="sm"
              className="w-full mt-2 rounded-lg border-gray-200 text-neutral-700 hover:border-gray-300 hover:text-neutral-900 transition-colors"
              onClick={() => navigate('/rules')}
            >
              View All Rules
            </Button>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5 space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold text-neutral-700">Recent Activity</h3>
            </div>
            {recentTransactions.length === 0 ? (
              <div className="text-center py-4 text-sm text-neutral-500">
                No recent transactions
              </div>
            ) : (
              <div className="space-y-2">
                {recentTransactions.map((transaction) => {
                  const amount = Number(transaction.amount) / 100
                  const isIncome = amount < 0 // Plaid returns negative amounts for credits
                  const displayAmount = Math.abs(amount).toFixed(2)

                  return (
                    <div key={transaction.id.toString()} className="flex items-center justify-between text-sm">
                      <div className="flex-1 min-w-0">
                        <div className="text-neutral-700 truncate">
                          {transaction.merchant_name || transaction.name}
                        </div>
                        <div className="text-xs text-neutral-500">
                          {new Date(transaction.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                          {transaction.category && ` • ${transaction.category}`}
                        </div>
                      </div>
                      <div className={`font-semibold ${isIncome ? 'text-green-600' : 'text-neutral-900'}`}>
                        {isIncome ? '+' : '-'}${displayAmount}
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </CardContent>
        </Card>
      </PageBody>
    </PageShell>
  )
}

function KPI({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-xl bg-gray-50 p-3 text-center border border-gray-100">
      <div className="text-neutral-500 text-xs">{label}</div>
      <div className="font-semibold text-neutral-900">{value}</div>
    </div>
  )
}

function Rule({ text, isActive = true }: { text: string; isActive?: boolean }) {
  return (
    <div className={`rounded-md bg-white px-3 py-2 text-sm border border-gray-100 hover:border-gray-200 transition-colors ${!isActive ? 'opacity-50' : ''}`}>
      {text}
    </div>
  )
}

export default Dashboard
