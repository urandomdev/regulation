import { useCallback, useEffect, useRef, useState } from 'react'
import { usePlaidLink } from 'react-plaid-link'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { SyncStatusModal, type SyncPhase, type LatestActivity } from '@/components/SyncStatusModal'

interface PlaidLinkOnSuccessMetadata {
  institution: null | {
    name: string
    institution_id: string
  }
  accounts: Array<{
    id: string
    name: string
    mask: string
    type: string
    subtype: string
    verification_status: string
  }>
  link_session_id: string
  transfer_status?: string
}

interface PlaidLinkOnExitMetadata {
  institution: null | {
    name: string
    institution_id: string
  }
  status: null | string
  link_session_id: string
  request_id: string
}

interface PlaidLinkError {
  error_type: string
  error_code: string
  error_message: string
  display_message: string
}

type PlaidLinkOnSuccess = (public_token: string, metadata: PlaidLinkOnSuccessMetadata) => void
type PlaidLinkOnExit = (error: null | PlaidLinkError, metadata: PlaidLinkOnExitMetadata) => void

interface PlaidLinkButtonProps {
  onSuccess?: () => void
  variant?: 'default' | 'outline' | 'ghost'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  className?: string
  children?: React.ReactNode
}

export function PlaidLinkButton({
  onSuccess: onSuccessCallback,
  variant = 'default',
  size = 'default',
  className,
  children = 'Connect Bank Account'
}: PlaidLinkButtonProps) {
  const [linkToken, setLinkToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [syncStatus, setSyncStatus] = useState<{
    open: boolean
    phase: SyncPhase
    message: string
    activities: LatestActivity[]
    syncedCount: number
  }>({ open: false, phase: 'idle', message: '', activities: [], syncedCount: 0 })
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const syncedIdsRef = useRef<Set<string>>(new Set())

  const stopPolling = useCallback(() => {
    if (pollingRef.current) {
      clearInterval(pollingRef.current)
      pollingRef.current = null
    }
  }, [])

  const closeSyncStatus = () => {
    stopPolling()
    syncedIdsRef.current = new Set()
    setSyncStatus((state) => ({ ...state, open: false, activities: [], syncedCount: 0, phase: 'idle', message: '' }))
  }

  const mergeActivities = useCallback((activities: LatestActivity[]) => {
    if (!activities.length) {
      return syncedIdsRef.current.size
    }
    const nextSet = new Set(syncedIdsRef.current)
    activities.forEach((activity) => nextSet.add(activity.id))
    syncedIdsRef.current = nextSet
    return nextSet.size
  }, [])

  const fetchLatestActivities = useCallback(async (): Promise<LatestActivity[]> => {
    const [data, err] = await api.financial.getTransactions({
      limit: BigInt(3),
      offset: BigInt(0),
    })

    if (err || !data?.transactions?.length) {
      return []
    }

    return data.transactions.slice(0, 3).map((tx) => ({
      id: tx.id.toString(),
      name: tx.merchant_name || tx.name,
      amountCents: Number(tx.amount),
      date: tx.date instanceof Date ? tx.date.toISOString() : tx.date?.toString?.() ?? new Date().toISOString(),
    }))
  }, [])

  // Callback when user successfully links account
  const onSuccess = useCallback<PlaidLinkOnSuccess>(
    async (public_token, metadata) => {
      console.log('Plaid Link Success:', metadata)

      try {
        // Exchange public token with backend
        const [, exchangeErr] = await api.plaid.exchangeToken({ public_token })

        if (exchangeErr) {
          console.error('Failed to exchange token:', exchangeErr)
          alert('Failed to connect account. Please try again.')
          return
        }

        syncedIdsRef.current = new Set()
        setSyncStatus({
          open: true,
          phase: 'syncing',
          message: 'We are pulling the latest transactions from your institution.',
          activities: [],
          syncedCount: 0,
        })

        // Trigger transaction sync for the newly connected account
        console.log('Triggering transaction sync...')
        const [, syncErr] = await api.plaid.syncTransactions({})

        if (syncErr) {
          console.error('Failed to sync transactions:', syncErr)
          stopPolling()
          setSyncStatus({
            open: true,
            phase: 'error',
            message: 'We connected your account but could not finish syncing transactions.',
            activities: [],
            syncedCount: syncedIdsRef.current.size,
          })
        } else {
          const latestActivities = await fetchLatestActivities()
          stopPolling()
          const mergedCount = mergeActivities(latestActivities)
          setSyncStatus({
            open: true,
            phase: 'success',
            message: latestActivities.length
              ? 'Your most recent activity is ready below.'
              : 'Sync completed. New activity will show up automatically.',
            activities: latestActivities,
            syncedCount: mergedCount,
          })
        }

        // Call success callback to refresh data
        if (onSuccessCallback) {
          onSuccessCallback()
        }
      } catch (error) {
        console.error('Error exchanging token:', error)
        alert('Failed to connect account. Please try again.')
      }
    },
    [fetchLatestActivities, mergeActivities, onSuccessCallback, stopPolling]
  )

  // Callback when user exits without linking
  const onExit = useCallback<PlaidLinkOnExit>(
    (error, metadata) => {
      if (error) {
        console.error('Plaid Link Error:', error)
      }
      console.log('Plaid Link Exit:', metadata)
      setLoading(false)
    },
    []
  )

  // Initialize Plaid Link
  const config = {
    token: linkToken,
    onSuccess,
    onExit,
  }

  const { open, ready } = usePlaidLink(config)

  // Auto-open Link when it becomes ready
  useEffect(() => {
    if (linkToken && ready && loading) {
      open()
      setLoading(false)
    }
  }, [linkToken, ready, loading, open])

  // Handle button click - create link token
  const handleClick = async () => {
    setLoading(true)

    try {
      // Create link token from backend
      const [data, err] = await api.plaid.createLinkToken()

      if (err || !data) {
        console.error('Failed to create link token:', err)
        alert('Failed to initialize bank connection. Please try again.')
        setLoading(false)
        return
      }

      setLinkToken(data.link_token)
      // Link will auto-open via useEffect when ready becomes true
    } catch (error) {
      console.error('Error creating link token:', error)
      alert('Failed to initialize bank connection. Please try again.')
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!syncStatus.open || syncStatus.phase !== 'syncing') {
      stopPolling()
      return
    }

    stopPolling()
    pollingRef.current = setInterval(async () => {
      const activities = await fetchLatestActivities()
      if (activities.length) {
        const mergedCount = mergeActivities(activities)
        setSyncStatus((state) =>
          state.phase === 'syncing' ? { ...state, activities, syncedCount: mergedCount } : state
        )
      }
    }, 3000)

    return () => {
      stopPolling()
    }
  }, [fetchLatestActivities, mergeActivities, stopPolling, syncStatus.open, syncStatus.phase])


  return (
    <>
      <Button
        onClick={handleClick}
        disabled={loading}
        variant={variant}
        size={size}
        className={className}
      >
        {loading ? 'Initializing...' : children}
      </Button>

      <SyncStatusModal
        open={syncStatus.open}
        phase={syncStatus.phase}
        message={syncStatus.message}
        activities={syncStatus.activities}
        totalSynced={syncStatus.syncedCount}
        onClose={closeSyncStatus}
      />
    </>
  )
}
