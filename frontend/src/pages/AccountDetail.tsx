import { useCallback, useEffect, useState } from 'react'
import { UUID } from '@deltalaboratory/uuid'
import { useParams } from 'react-router'
import { TrendingUp, TrendingDown } from 'lucide-react'
import { api } from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
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

interface Transaction {
  id: UUID
  account_id: UUID
  amount: bigint
  date: string
  name: string
  merchant_name?: string
  category: string
  pending: boolean
  payment_channel?: string
}

const PAGE_SIZE = 20

const AccountDetail = () => {
  const { accountId } = useParams<{ accountId: string }>()
  const [account, setAccount] = useState<Account | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [transactionsError, setTransactionsError] = useState<string | null>(null)
  const [transactionsLoading, setTransactionsLoading] = useState(false)
  const [page, setPage] = useState(0)
  const [hasMore, setHasMore] = useState(false)

  const fetchTransactions = useCallback(async (pageIndex: number) => {
    if (!accountId) return

    setTransactionsLoading(true)
    setTransactionsError(null)

    const [transactionsData, transactionsErr] = await api.financial.getTransactions({
      account_id: new UUID(accountId),
      limit: BigInt(PAGE_SIZE),
      offset: BigInt(pageIndex * PAGE_SIZE),
    })

    if (transactionsErr) {
      setTransactions([])
      setHasMore(false)
      setTransactionsError(transactionsErr.message || 'Failed to fetch transactions')
      setTransactionsLoading(false)
      return
    }

    const convertedTransactions: Transaction[] = (transactionsData?.transactions || []).map(tx => ({
      id: tx.id,
      account_id: tx.account_id,
      amount: tx.amount,
      date: tx.date instanceof Date ? tx.date.toISOString() : new Date(tx.date).toISOString(),
      name: tx.name,
      merchant_name: tx.merchant_name ?? undefined,
      category: tx.category,
      pending: tx.pending,
      payment_channel: tx.payment_channel ?? undefined
    }))

    setTransactions(convertedTransactions)
    setHasMore((transactionsData?.transactions?.length || 0) === PAGE_SIZE)
    setTransactionsLoading(false)
  }, [accountId])

  useEffect(() => {
    const fetchAccountData = async () => {
      if (!accountId) {
        setError('Account ID is required')
        setLoading(false)
        return
      }

      setLoading(true)

      // Fetch account details
      const [accountsData, accountsErr] = await api.financial.getAccounts()
      if (accountsErr) {
        setError(accountsErr.message || 'Failed to fetch account')
        setLoading(false)
        return
      }

      if (accountsData) {
        const foundAccount = accountsData.accounts.find(
          (acc: Account) => acc.id.toString() === accountId
        )
        if (foundAccount) {
          setAccount(foundAccount)
        } else {
          setError('Account not found')
          setLoading(false)
          return
        }
      }

      setLoading(false)
    }

    fetchAccountData()
  }, [accountId])

  useEffect(() => {
    if (!accountId || !account) {
      return
    }
    fetchTransactions(page)
  }, [account, accountId, fetchTransactions, page])

  useEffect(() => {
    setPage(0)
    setTransactions([])
  }, [accountId])

  const formatBalance = (balance: bigint) => {
    return (Number(balance) / 100).toFixed(2)
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  }

  const getAccountTypeLabel = (type: string) => {
    return type.charAt(0).toUpperCase() + type.slice(1).toLowerCase()
  }

  const handleNextPage = () => {
    if (!hasMore || transactionsLoading) return
    setPage((prev) => prev + 1)
  }

  const handlePrevPage = () => {
    if (page === 0 || transactionsLoading) return
    setPage((prev) => Math.max(prev - 1, 0))
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

  if (!account) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <p className="text-neutral-500">Account not found</p>
        </PageBody>
      </PageShell>
    )
  }

  return (
    <PageShell>
      <PageHeader title="Account Details" backTo="/accounts" />

      <PageBody className="space-y-4">
        {/* Account Summary Card */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-6">
            <div className="text-center">
              <div className="text-sm text-neutral-500 mb-1">{account.name}</div>
              <div className="text-xs text-neutral-400 mb-4">
                {getAccountTypeLabel(account.type)}
                {account.mask && ` ••${account.mask}`}
              </div>
              <div className="text-4xl font-bold text-neutral-900 mb-2">
                ${formatBalance(account.current_balance)}
              </div>
              {account.available_balance !== null && account.available_balance !== undefined && (
                <div className="text-sm text-neutral-500">
                  ${formatBalance(account.available_balance)} available
                </div>
              )}
            </div>
            {!account.is_active && (
              <div className="mt-4 text-sm text-amber-600 bg-amber-50 px-3 py-2 rounded-lg text-center">
                This account is currently inactive
              </div>
            )}
          </CardContent>
        </Card>

        {/* Transactions Section */}
        <Card className="rounded-2xl shadow-sm border border-gray-100">
          <CardContent className="p-5 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-sm font-semibold text-neutral-700">Transactions</h3>
                <p className="text-xs text-neutral-500">Showing {PAGE_SIZE} per page</p>
              </div>
              {transactionsLoading && (
                <span className="text-xs text-neutral-400">Refreshing…</span>
              )}
            </div>

            {transactionsError ? (
              <div className="rounded-lg border border-red-100 bg-red-50 px-3 py-2 text-sm text-red-600">
                {transactionsError}
              </div>
            ) : transactionsLoading && transactions.length === 0 ? (
              <div className="flex items-center justify-center py-8 text-sm text-neutral-500">
                Loading transactions…
              </div>
            ) : transactions.length === 0 ? (
              <div className="text-center py-8 text-neutral-500">
                No transactions found
              </div>
            ) : (
              <>
                <div className="space-y-3">
                  {transactions.map((transaction) => {
                    const amount = Number(transaction.amount) / 100
                    const isIncome = amount < 0 // Plaid returns negative amounts for credits

                    return (
                      <div
                        key={transaction.id.toString()}
                        className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0"
                      >
                        <div className="flex items-center gap-3 flex-1">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                            isIncome ? 'bg-green-50' : 'bg-gray-50'
                          }`}>
                            {isIncome ? (
                              <TrendingUp className="h-4 w-4 text-green-600" />
                            ) : (
                              <TrendingDown className="h-4 w-4 text-gray-500" />
                            )}
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="font-medium text-sm text-neutral-900 truncate">
                              {transaction.merchant_name || transaction.name}
                            </div>
                            <div className="text-xs text-neutral-500">
                              {formatDate(transaction.date)}
                              {transaction.pending && (
                                <span className="ml-2 text-amber-600">• Pending</span>
                              )}
                            </div>
                            <div className="text-xs text-neutral-400 mt-0.5">
                              {transaction.category}
                            </div>
                          </div>
                        </div>
                        <div className={`font-semibold text-sm ${
                          isIncome ? 'text-green-600' : 'text-neutral-900'
                        }`}>
                          {isIncome ? '+' : '-'}${Math.abs(amount).toFixed(2)}
                        </div>
                      </div>
                    )
                  })}
                </div>

                <div className="border-t border-gray-100 pt-3">
                  <div className="flex items-center justify-between text-xs text-neutral-500">
                    <span>Page {page + 1}</span>
                    <span>
                      {hasMore ? 'More transactions available' : 'End of history'}
                    </span>
                  </div>
                  <div className="mt-3 flex items-center justify-between gap-3">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handlePrevPage}
                      disabled={page === 0 || transactionsLoading}
                    >
                      Previous
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleNextPage}
                      disabled={!hasMore || transactionsLoading}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </PageBody>
    </PageShell>
  )
}

export default AccountDetail
