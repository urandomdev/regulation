import { useEffect, useState } from 'react'
import { UUID } from '@deltalaboratory/uuid'
import { Plus, Trash2, Power, Edit2 } from 'lucide-react'
import { api } from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { PageShell, PageHeader, PageBody, EmptyState } from '@/components/layout/Page'
import type { Category, ActionType } from '../../../api/sdk/js/dist/esm/regulation/internal/ent/rule'

// Helper to convert string to SDK Category type
const categoryToSDK = (category: string): Category => {
  return category.charAt(0).toUpperCase() + category.slice(1) as Category
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

interface Account {
  id: UUID
  name: string
  type: string
}

const Rules = () => {
  const [rules, setRules] = useState<Rule[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingRule, setEditingRule] = useState<Rule | null>(null)

  const fetchData = async () => {
    setLoading(true)

    // Fetch rules
    const [rulesData, rulesErr] = await api.rule.listRules()
    if (!rulesErr && rulesData) {
      setRules(rulesData.rules)
    }

    // Fetch accounts for the form
    const [accountsData, accountsErr] = await api.financial.getAccounts()
    if (!accountsErr && accountsData) {
      setAccounts(accountsData.accounts)
    }

    if (rulesErr) {
      setError(rulesErr.message || 'Failed to fetch rules')
    }

    setLoading(false)
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleToggle = async (ruleId: UUID) => {
    const [, err] = await api.rule.toggleRule(ruleId.toString())
    if (err) {
      alert(`Failed to toggle rule: ${err.message || 'Unknown error'}`)
      return
    }
    fetchData()
  }

  const handleDelete = async (rule: Rule) => {
    if (!confirm(`Are you sure you want to delete "${rule.name}"?`)) {
      return
    }

    const [, err] = await api.rule.deleteRule(rule.id.toString())
    if (err) {
      alert(`Failed to delete rule: ${err.message || 'Unknown error'}`)
      return
    }
    fetchData()
  }

  const formatCurrency = (cents: bigint) => {
    return (Number(cents) / 100).toFixed(2)
  }

  const getCategoryLabel = (category: string) => {
    return category.charAt(0).toUpperCase() + category.slice(1).toLowerCase()
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

  const savingsAccounts = accounts.filter((a) => a.type === 'savings')

  const activeCount = rules.filter((rule) => rule.is_active).length
  const totalSaved = rules.reduce((sum, rule) => sum + Number(rule.total_saved_cents ?? 0n), 0)
  const averageSaved = rules.length ? totalSaved / rules.length : 0

  if (showForm) {
    return (
      <RuleForm
        accounts={savingsAccounts}
        rule={editingRule}
        onCancel={() => {
          setShowForm(false)
          setEditingRule(null)
        }}
        onSuccess={() => {
          setShowForm(false)
          setEditingRule(null)
          fetchData()
        }}
      />
    )
  }

  return (
    <PageShell>
      <PageHeader
        title="Savings Rules"
        backTo="/dashboard"
        action={(
          <Button size="sm" onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-1" />
            Add Rule
          </Button>
        )}
      />

      <PageBody className="space-y-4">
        <Card className="rounded-2xl border border-gray-100 shadow-sm">
          <CardContent className="grid grid-cols-3 gap-3 p-4 text-center text-sm">
            <div>
              <p className="text-xs text-neutral-500 uppercase tracking-wide">Active Rules</p>
              <p className="text-lg font-semibold text-neutral-900">{activeCount}</p>
            </div>
            <div>
              <p className="text-xs text-neutral-500 uppercase tracking-wide">Total Saved</p>
              <p className="text-lg font-semibold text-neutral-900">${(totalSaved / 100).toFixed(2)}</p>
            </div>
            <div>
              <p className="text-xs text-neutral-500 uppercase tracking-wide">Avg / Rule</p>
              <p className="text-lg font-semibold text-neutral-900">${averageSaved ? (averageSaved / 100).toFixed(2) : '0.00'}</p>
            </div>
          </CardContent>
        </Card>

        {rules.length === 0 ? (
          <EmptyState
            title="No rules yet"
            description="Create your first automation to move money every time you spend."
            action={(
              <Button onClick={() => setShowForm(true)}>
                Create Your First Rule
              </Button>
            )}
          />
        ) : (
          <div className="space-y-3">
            {rules.map((rule) => (
              <Card
                key={rule.id.toString()}
                className={`rounded-2xl border shadow-sm ${rule.is_active ? 'border-gray-100' : 'border-gray-200 bg-gray-50'}`}
              >
                <CardContent className="p-4 space-y-3">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="flex items-center gap-2">
                        <h3 className={`text-base font-semibold ${rule.is_active ? 'text-neutral-900' : 'text-neutral-500'}`}>
                          {rule.name}
                        </h3>
                        <span className="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-neutral-600">
                          {getCategoryLabel(rule.category)}
                        </span>
                        {!rule.is_active && (
                          <span className="text-xs px-2 py-0.5 bg-gray-200 text-gray-600 rounded-full">Inactive</span>
                        )}
                      </div>
                      <p className="text-xs text-neutral-500 mt-1">
                        {rule.min_amount_cents && `Min $${formatCurrency(rule.min_amount_cents)} `}
                        {rule.max_amount_cents && `Max $${formatCurrency(rule.max_amount_cents)}`}
                      </p>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className={`transition-colors ${rule.is_active ? 'text-neutral-900' : 'text-neutral-400'}`}
                      onClick={() => handleToggle(rule.id)}
                      aria-label={rule.is_active ? 'Deactivate rule' : 'Activate rule'}
                    >
                      <Power className="h-4 w-4" />
                    </Button>
                  </div>

                  <div className="rounded-lg bg-gray-50 px-3 py-2 text-xs text-neutral-600">
                    <div className="flex items-center justify-between">
                      <span>
                        {rule.action_type === 'multiply'
                          ? `Multiply by ${rule.action_value}x`
                          : `Save $${rule.action_value.toFixed(2)}`}
                      </span>
                      <span>Priority {Number(rule.priority)}</span>
                    </div>
                    <div className="mt-1 text-[11px] text-neutral-500">
                      Executed {Number(rule.execution_count)} times • Saved ${formatCurrency(rule.total_saved_cents)} total
                    </div>
                  </div>

                  <div className="flex gap-2 pt-2">
                    <Button
                      variant="outline"
                      size="sm"
                      className="flex-1"
                      onClick={() => {
                        setEditingRule(rule)
                        setShowForm(true)
                      }}
                    >
                      <Edit2 className="h-4 w-4 mr-1" />
                      Edit
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="flex-1 border-red-200 text-red-600 hover:bg-red-50"
                      onClick={() => handleDelete(rule)}
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Delete
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </PageBody>
    </PageShell>
  )
}

interface RuleFormProps {
  accounts: Account[]
  rule: Rule | null
  onCancel: () => void
  onSuccess: () => void
}

function RuleForm({ accounts, rule, onCancel, onSuccess }: RuleFormProps) {
  const [name, setName] = useState(rule?.name || '')
  const [category, setCategory] = useState(rule?.category || 'dining')
  const [minAmount, setMinAmount] = useState(rule?.min_amount_cents ? Number(rule.min_amount_cents) / 100 : '')
  const [maxAmount, setMaxAmount] = useState(rule?.max_amount_cents ? Number(rule.max_amount_cents) / 100 : '')
  const [actionType, setActionType] = useState(rule?.action_type || 'multiply')
  const [actionValue, setActionValue] = useState(rule?.action_value || 1)
  const [targetAccountId, setTargetAccountId] = useState(rule?.target_account_id.toString() || '')
  const [priority, setPriority] = useState(rule?.priority ? Number(rule.priority) : 0)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    const params = {
      name,
      category: categoryToSDK(category),
      min_amount_cents: minAmount ? BigInt(Math.round(Number(minAmount) * 100)) : undefined,
      max_amount_cents: maxAmount ? BigInt(Math.round(Number(maxAmount) * 100)) : undefined,
      action_type: actionType as ActionType,
      action_value: Number(actionValue),
      target_account_id: new UUID(targetAccountId),
      priority: BigInt(priority),
    }

    if (rule) {
      const [, err] = await api.rule.updateRule(rule.id.toString(), params)
      if (err) {
        alert(`Failed to update rule: ${err.message || 'Unknown error'}`)
        setSubmitting(false)
        return
      }
    } else {
      const [, err] = await api.rule.createRule(params)
      if (err) {
        alert(`Failed to create rule: ${err.message || 'Unknown error'}`)
        setSubmitting(false)
        return
      }
    }

    onSuccess()
  }

  const categories = ['dining', 'groceries', 'transport', 'shopping', 'subscriptions', 'entertainment', 'bills', 'misc']

  const selectClassName = 'h-10 w-full rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2'

  return (
    <PageShell>
      <PageHeader title={rule ? 'Edit Rule' : 'Create Rule'} backTo={onCancel} />

      <PageBody>
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5">
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="rule-name">Rule Name</Label>
                <Input
                  id="rule-name"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g., Coffee roundup"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="rule-category">Category</Label>
                <select
                  id="rule-category"
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  required
                  className={selectClassName}
                >
                  {categories.map((cat) => (
                    <option key={cat} value={cat}>
                      {cat.charAt(0).toUpperCase() + cat.slice(1)}
                    </option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="min-amount">Min Amount ($)</Label>
                  <Input
                    id="min-amount"
                    type="number"
                    step="0.01"
                    value={minAmount}
                    onChange={(e) => setMinAmount(e.target.value)}
                    placeholder="Optional"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="max-amount">Max Amount ($)</Label>
                  <Input
                    id="max-amount"
                    type="number"
                    step="0.01"
                    value={maxAmount}
                    onChange={(e) => setMaxAmount(e.target.value)}
                    placeholder="Optional"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="action-type">Action Type</Label>
                <select
                  id="action-type"
                  value={actionType}
                  onChange={(e) => setActionType(e.target.value)}
                  required
                  className={selectClassName}
                >
                  <option value="multiply">Multiply (e.g., roundup)</option>
                  <option value="fixed">Fixed amount</option>
                </select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="action-value">{actionType === 'multiply' ? 'Multiplier' : 'Amount ($)'}</Label>
                <Input
                  id="action-value"
                  type="number"
                  step={actionType === 'multiply' ? '0.1' : '0.01'}
                  value={actionValue}
                  onChange={(e) => setActionValue(Number(e.target.value))}
                  placeholder={actionType === 'multiply' ? 'e.g., 2' : 'e.g., 5.00'}
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="target-account">Target Savings Account</Label>
                <select
                  id="target-account"
                  value={targetAccountId}
                  onChange={(e) => setTargetAccountId(e.target.value)}
                  required
                  className={selectClassName}
                >
                  <option value="">Select account</option>
                  {accounts.map((account) => (
                    <option key={account.id.toString()} value={account.id.toString()}>
                      {account.name}
                    </option>
                  ))}
                </select>
                {accounts.length === 0 && (
                  <p className="text-xs text-amber-600">No savings accounts found. Please add one first.</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="priority">Priority</Label>
                <Input
                  id="priority"
                  type="number"
                  value={priority}
                  onChange={(e) => setPriority(Number(e.target.value))}
                  placeholder="0 (higher = more priority)"
                  required
                />
              </div>

              <div className="flex flex-col gap-3 pt-4 sm:flex-row">
                <Button
                  type="button"
                  variant="outline"
                  className="flex-1"
                  onClick={onCancel}
                  disabled={submitting}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  className="flex-1"
                  disabled={submitting || accounts.length === 0}
                >
                  {submitting ? 'Saving...' : rule ? 'Update Rule' : 'Create Rule'}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </PageBody>
    </PageShell>
  )
}

export default Rules
