import { useEffect, useState } from 'react'
import { UUID } from '@deltalaboratory/uuid'
import { useNavigate } from 'react-router'
import { CreditCard, Plus, ChevronRight } from 'lucide-react'
import { api } from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { PlaidLinkButton } from '@/components/PlaidLink'
import { PageShell, PageHeader, PageBody } from '@/components/layout/Page'

interface Account {
  id: UUID
  name: string
  type: string
  mask?: string | null
  current_balance: bigint
  available_balance?: bigint | null
  is_active: boolean
}

const Accounts = () => {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  const fetchAccounts = async () => {
    setLoading(true)
    const [data, err] = await api.financial.getAccounts()
    if (err) {
      setError(err.message || 'Failed to fetch accounts')
    } else if (data) {
      setAccounts(data.accounts)
    }
    setLoading(false)
  }

  useEffect(() => {
    fetchAccounts()
  }, [])

  const formatBalance = (balance: bigint) => {
    return (Number(balance) / 100).toFixed(2)
  }

  const getAccountTypeLabel = (type: string) => {
    return type.charAt(0).toUpperCase() + type.slice(1).toLowerCase()
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

  return (
    <PageShell>
      <PageHeader
        title="Accounts"
        backTo="/dashboard"
        action={(
          <PlaidLinkButton
            variant="ghost"
            size="sm"
            className="text-sm text-neutral-500 hover:text-neutral-700 transition-colors"
            onSuccess={fetchAccounts}
          >
            <Plus className="h-4 w-4 mr-1" />
            Add
          </PlaidLinkButton>
        )}
      />

      <PageBody>
        {accounts.length === 0 ? (
          <Card className="rounded-2xl shadow-sm border border-gray-100">
            <CardContent className="p-8 text-center">
              <CreditCard className="h-12 w-12 mx-auto text-gray-300 mb-3" />
              <p className="text-neutral-600 mb-4">No accounts connected yet</p>
              <PlaidLinkButton
                onSuccess={fetchAccounts}
                className="rounded-lg"
              >
                Connect Your First Account
              </PlaidLinkButton>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {accounts.map((account) => (
              <Card
                key={account.id.toString()}
                className="rounded-2xl shadow-sm border border-gray-100 hover:border-gray-200 transition-colors"
              >
                <CardContent className="p-4">
                  <div
                    className="flex items-center justify-between cursor-pointer"
                    onClick={() => navigate(`/accounts/${account.id.toString()}`)}
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center">
                        <CreditCard className="h-5 w-5 text-gray-500" />
                      </div>
                      <div>
                        <div className="font-semibold text-neutral-900">{account.name}</div>
                        <div className="text-sm text-neutral-500">
                          {getAccountTypeLabel(account.type)}
                          {account.mask && ` ••${account.mask}`}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="text-right">
                        <div className="font-semibold text-neutral-900">
                          ${formatBalance(account.current_balance)}
                        </div>
                        {account.available_balance !== null && account.available_balance !== undefined && (
                          <div className="text-xs text-neutral-500">
                            ${formatBalance(account.available_balance)} available
                          </div>
                        )}
                      </div>
                      <ChevronRight className="h-5 w-5 text-neutral-300" />
                    </div>
                  </div>
                  {!account.is_active && (
                    <div className="mt-2 text-xs text-amber-600 bg-amber-50 px-2 py-1 rounded">
                      Inactive
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Total Summary */}
        {accounts.length > 0 && (
          <Card className="rounded-2xl shadow-sm border border-gray-100 mt-6">
            <CardContent className="p-4">
              <div className="flex justify-between items-center">
                <span className="text-neutral-600">Total Balance</span>
                <span className="text-xl font-bold text-neutral-900">
                  ${formatBalance(accounts.reduce((sum, acc) => sum + acc.current_balance, 0n))}
                </span>
              </div>
            </CardContent>
          </Card>
        )}

        <div className="mt-6 text-center text-xs text-neutral-500">
          Need to remove a connection?{' '}
          <button
            type="button"
            className="text-neutral-900 underline-offset-4 hover:underline"
            onClick={() => navigate('/accounts/disconnect')}
          >
            Manage disconnects
          </button>
        </div>
      </PageBody>
    </PageShell>
  )
}

export default Accounts
