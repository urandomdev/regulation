import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router'
import { UUID } from '@deltalaboratory/uuid'

import { PageShell, PageHeader, PageBody, EmptyState } from '@/components/layout/Page'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'
import type { Category as SDKCategory, ActionType as SDKActionType } from '../../../api/sdk/js/dist/esm/regulation/internal/ent/rule'

interface RuleSuggestionResponse {
  name: string
  category: string
  action_type: string
  action_value: number
  min_amount_cents?: number | null
  max_amount_cents?: number | null
  estimated_savings: number
  confidence: string
  reasoning: string
  impact_level: string
}

type NumericInput = number | bigint | null | undefined

interface AccountOption {
  id: string
  name: string
}

const toOptionalNumber = (value: NumericInput): number | null => {
  if (value === null || value === undefined) {
    return null
  }
  return typeof value === 'bigint' ? Number(value) : value
}

const toNumberOr = (value: NumericInput, fallback = 0): number => {
  const result = toOptionalNumber(value)
  return result === null ? fallback : result
}

const toSDKCategory = (category: string): SDKCategory => {
  if (!category) {
    return 'Misc' as SDKCategory
  }
  const normalized = category.charAt(0).toUpperCase() + category.slice(1).toLowerCase()
  return normalized as SDKCategory
}

const toSDKActionType = (actionType: string): SDKActionType => {
  if (actionType !== 'multiply' && actionType !== 'fixed') {
    return 'multiply' as SDKActionType
  }
  return actionType as SDKActionType
}

const Recommendations = () => {
  const navigate = useNavigate()
  const [suggestions, setSuggestions] = useState<RuleSuggestionResponse[]>([])
  const [analysisSummary, setAnalysisSummary] = useState<string>('')
  const [prioritySuggestion, setPrioritySuggestion] = useState<string>('')
  const [analysisPeriod, setAnalysisPeriod] = useState<number>(30)
  const [transactionCount, setTransactionCount] = useState<number>(0)
  const [totalSpent, setTotalSpent] = useState<number>(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [accountsError, setAccountsError] = useState<string | null>(null)
  const [savingsAccounts, setSavingsAccounts] = useState<AccountOption[]>([])
  const [selectedAccounts, setSelectedAccounts] = useState<Record<number, string>>({})
  const [creatingIndex, setCreatingIndex] = useState<number | null>(null)
  const [completed, setCompleted] = useState<Record<number, boolean>>({})
  const [creationErrors, setCreationErrors] = useState<Record<number, string>>({})

  const fetchRecommendations = async () => {
    setLoading(true)
    setError(null)
    setAccountsError(null)

    const [data, err] = await api.recommendation.getRecommendations()
    const [accountsData, accountsErr] = await api.financial.getAccounts()

    if (err || !data) {
      setSuggestions([])
      setAnalysisSummary('')
      setPrioritySuggestion('')
      setError(err?.message || 'Need at least five settled transactions to generate recommendations.')
      setLoading(false)
      return
    }

    let nextSavings: AccountOption[] = []
    if (accountsErr || !accountsData) {
      setAccountsError(accountsErr?.message || 'Unable to load savings accounts. Add a savings account to build rules here.')
      setSavingsAccounts([])
    } else {
      nextSavings = (accountsData.accounts || [])
        .filter((account) => account.type === 'savings')
        .map((account) => ({
          id: account.id.toString(),
          name: account.name,
        }))

      setSavingsAccounts(nextSavings)
    }

    const normalizedSuggestions: RuleSuggestionResponse[] = (data.suggestions ?? []).map((suggestion) => ({
      ...suggestion,
      min_amount_cents: toOptionalNumber(suggestion.min_amount_cents),
      max_amount_cents: toOptionalNumber(suggestion.max_amount_cents),
      action_value: toNumberOr(suggestion.action_value),
      estimated_savings: toNumberOr(suggestion.estimated_savings),
    }))

    setSuggestions(normalizedSuggestions)
    setAnalysisSummary(data.overall_analysis || '')
    setPrioritySuggestion(data.priority_suggestion || '')
    setAnalysisPeriod(toNumberOr(data.analysis_period_days, 30))
    setTransactionCount(toNumberOr(data.transaction_count))
    setTotalSpent(toNumberOr(data.total_spent))
    const defaultAccountId = nextSavings[0]?.id
    if (defaultAccountId) {
      setSelectedAccounts((prev) => {
        const nextSelections: Record<number, string> = {}
        normalizedSuggestions.forEach((_, index) => {
          nextSelections[index] = prev[index] || defaultAccountId
        })
        return nextSelections
      })
    }
    setCompleted({})
    setCreationErrors({})
    setLoading(false)
  }
  const defaultAccountId = useMemo(() => savingsAccounts[0]?.id ?? '', [savingsAccounts])

  const handleAccountChange = (index: number, accountId: string) => {
    setSelectedAccounts((prev) => ({
      ...prev,
      [index]: accountId,
    }))
  }

  const handleCreateRule = async (rec: RuleSuggestionResponse, index: number) => {
    const accountId = selectedAccounts[index] || defaultAccountId

    if (!accountId) {
      setCreationErrors((prev) => ({
        ...prev,
        [index]: 'Add a savings account before building this rule.',
      }))
      return
    }

    setCreatingIndex(index)
    setCreationErrors((prev) => {
      const next = { ...prev }
      delete next[index]
      return next
    })

    const params = {
      name: rec.name,
      category: toSDKCategory(rec.category),
      min_amount_cents: rec.min_amount_cents != null ? BigInt(Math.round(rec.min_amount_cents)) : undefined,
      max_amount_cents: rec.max_amount_cents != null ? BigInt(Math.round(rec.max_amount_cents)) : undefined,
      action_type: toSDKActionType(rec.action_type),
      action_value: Number(rec.action_value),
      target_account_id: new UUID(accountId),
      priority: BigInt(0),
    }

    const [, err] = await api.rule.createRule(params)

    if (err) {
      setCreationErrors((prev) => ({
        ...prev,
        [index]: err.message || 'Failed to create rule',
      }))
    } else {
      setCompleted((prev) => ({
        ...prev,
        [index]: true,
      }))
    }

    setCreatingIndex(null)
  }

  useEffect(() => {
    fetchRecommendations()
  }, [])

  return (
    <PageShell>
      <PageHeader title="Personalized rule ideas" backTo="/dashboard" />

      <PageBody className="space-y-4">
        <Card className="rounded-2xl border border-gray-100 shadow-sm">
          <CardContent className="grid grid-cols-2 gap-3 p-4 text-sm text-neutral-600">
            <div>
              <p className="text-xs uppercase tracking-wide text-neutral-500">Analysis window</p>
              <p className="text-lg font-semibold text-neutral-900">Last {analysisPeriod} days</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-neutral-500">Transactions reviewed</p>
              <p className="text-lg font-semibold text-neutral-900">{transactionCount}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-neutral-500">Total spend</p>
              <p className="text-lg font-semibold text-neutral-900">${totalSpent.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-neutral-500">Suggestions ready</p>
              <p className="text-lg font-semibold text-neutral-900">{suggestions.length}</p>
            </div>
          </CardContent>
        </Card>

        {accountsError && (
          <Card className="rounded-2xl border border-amber-200 bg-amber-50">
            <CardContent className="p-4 text-sm text-amber-700">
              {accountsError}
            </CardContent>
          </Card>
        )}

        {loading ? (
          <Card className="rounded-2xl border border-gray-100">
            <CardContent className="p-6 text-sm text-neutral-500">
              Analyzing your recent transactions…
            </CardContent>
          </Card>
        ) : error ? (
          <EmptyState
            title="Not enough data"
            description={error}
            action={(
              <Button onClick={fetchRecommendations}>Try again</Button>
            )}
          />
        ) : suggestions.length === 0 ? (
          <EmptyState
            title="No automation ideas yet"
            description="Keep transacting for a few more days and we’ll surface savings rules tailored to your habits."
          />
        ) : (
          <div className="space-y-4">
            {prioritySuggestion && (
              <Card className="rounded-2xl border border-gray-100 bg-gray-50">
                <CardContent className="p-4">
                  <p className="text-xs uppercase tracking-wide text-neutral-500">Priority suggestion</p>
                  <p className="mt-1 text-sm text-neutral-700">{prioritySuggestion}</p>
                </CardContent>
              </Card>
            )}

            {analysisSummary && (
              <Card className="rounded-2xl border border-gray-100">
                <CardContent className="p-4 text-sm text-neutral-600">
                  {analysisSummary}
                </CardContent>
              </Card>
            )}

            {suggestions.map((rec, index) => (
              <Card key={`${rec.name}-${index}`} className="rounded-2xl border border-gray-100">
                <CardContent className="space-y-2 p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-base font-semibold text-neutral-900">{rec.name}</p>
                      <p className="text-xs text-neutral-500">{rec.category} • {rec.impact_level}</p>
                    </div>
                    <span className="text-xs text-neutral-500 capitalize">{rec.confidence} confidence</span>
                  </div>

                  <p className="text-sm text-neutral-700">
                      {rec.action_type === 'multiply'
                        ? `Multiply by ${rec.action_value}x`
                        : `Save $${rec.action_value.toFixed(2)}`}
                    {rec.min_amount_cents != null && ` • Min $${(rec.min_amount_cents / 100).toFixed(2)}`}
                    {rec.max_amount_cents != null && ` • Max $${(rec.max_amount_cents / 100).toFixed(2)}`}
                  </p>

                  <p className="text-xs text-neutral-500">Estimated monthly savings ${rec.estimated_savings.toFixed(2)}</p>
                  <p className="text-xs text-neutral-500">{rec.reasoning}</p>

                  <div className="space-y-2 pt-2">
                    {savingsAccounts.length > 0 && (
                      <div className="text-xs text-neutral-600">
                        <label className="block font-medium text-neutral-500" htmlFor={`account-${index}`}>
                          Target savings account
                        </label>
                        <select
                          id={`account-${index}`}
                          className="mt-1 w-full rounded-lg border border-gray-200 bg-white px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-neutral-300"
                          value={selectedAccounts[index] || defaultAccountId}
                          onChange={(event) => handleAccountChange(index, event.target.value)}
                          disabled={!!completed[index]}
                        >
                          {savingsAccounts.map((account) => (
                            <option key={account.id} value={account.id}>
                              {account.name}
                            </option>
                          ))}
                        </select>
                      </div>
                    )}

                    {creationErrors[index] && (
                      <p className="text-xs text-red-600">{creationErrors[index]}</p>
                    )}

                    {completed[index] && !creationErrors[index] && (
                      <p className="text-xs text-emerald-600">Rule created successfully.</p>
                    )}

                    <Button
                      variant="outline"
                      size="sm"
                      className="w-full"
                      disabled={creatingIndex === index || completed[index]}
                      onClick={() => handleCreateRule(rec, index)}
                    >
                      {completed[index] ? 'Rule created' : creatingIndex === index ? 'Creating…' : 'Build this rule'}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}

            <div className="text-center">
              <Button variant="ghost" onClick={() => navigate('/dashboard')}>
                Skip to dashboard
              </Button>
            </div>
          </div>
        )}
      </PageBody>
    </PageShell>
  )
}

export default Recommendations
