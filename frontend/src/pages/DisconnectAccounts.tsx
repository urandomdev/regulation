import { useEffect, useState } from 'react'
import { UUID } from '@deltalaboratory/uuid'

import { PageShell, PageHeader, PageBody, EmptyState } from '@/components/layout/Page'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'

interface Account {
  id: UUID
  name: string
  type: string
  mask?: string | null
}

const DisconnectAccounts = () => {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [disconnectingId, setDisconnectingId] = useState<string | null>(null)

  const fetchAccounts = async () => {
    setLoading(true)
    setError(null)
    const [data, err] = await api.financial.getAccounts()
    if (err) {
      setError(err.message || 'Unable to load accounts right now.')
      setAccounts([])
    } else {
      setAccounts(data?.accounts ?? [])
    }
    setLoading(false)
  }

  const handleDisconnect = async (account: Account) => {
    if (!confirm(`Disconnect ${account.name}? This removes the account and its transactions.`)) {
      return
    }
    setDisconnectingId(account.id.toString())
    const [, err] = await api.plaid.disconnectAccount(account.id.toString())
    if (err) {
      alert(err.message || 'Failed to disconnect. Try again later.')
      setDisconnectingId(null)
      return
    }
    await fetchAccounts()
    setDisconnectingId(null)
  }

  useEffect(() => {
    fetchAccounts()
  }, [])

  if (loading) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <p className="text-sm text-neutral-500">Loading accounts…</p>
        </PageBody>
      </PageShell>
    )
  }

  if (error) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <div className="text-center space-y-3">
            <p className="text-sm text-destructive">{error}</p>
            <Button size="sm" onClick={fetchAccounts}>Retry</Button>
          </div>
        </PageBody>
      </PageShell>
    )
  }

  return (
    <PageShell>
      <PageHeader title="Manage connections" backTo="/accounts" />

      <PageBody className="space-y-4">
        <Card className="rounded-2xl border border-gray-100">
          <CardContent className="p-5 text-sm text-neutral-600 space-y-2">
            <p>Disconnecting removes the financial account and any synced transactions from Regulation.</p>
            <p className="text-xs text-neutral-500">Only use this if you no longer want to automate activity from that institution.</p>
          </CardContent>
        </Card>

        {accounts.length === 0 ? (
          <EmptyState
            title="No accounts to disconnect"
            description="Connect a bank first from the Accounts page."
          />
        ) : (
          <div className="space-y-3">
            {accounts.map((account) => (
              <Card key={account.id.toString()} className="rounded-2xl border border-gray-100">
                <CardContent className="flex items-center justify-between gap-4 p-4">
                  <div>
                    <p className="font-medium text-neutral-900">{account.name}</p>
                    <p className="text-xs text-neutral-500">
                      {account.type} {account.mask ? `• •${account.mask}` : ''}
                    </p>
                  </div>
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => handleDisconnect(account)}
                    disabled={disconnectingId === account.id.toString()}
                  >
                    {disconnectingId === account.id.toString() ? 'Disconnecting…' : 'Disconnect'}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </PageBody>
    </PageShell>
  )
}

export default DisconnectAccounts
