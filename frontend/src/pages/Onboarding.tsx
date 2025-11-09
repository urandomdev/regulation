import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { ShieldCheck, CreditCard, Sparkles, CheckCircle2, Bell } from 'lucide-react'

import { PageShell, PageHeader, PageBody, EmptyState } from '@/components/layout/Page'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { PlaidLinkButton } from '@/components/PlaidLink'
import { subscribeToNotifications } from '@/lib/notifications'
import { useAuth } from '@/hooks/useAuth'
import { api } from '@/lib/api'

const steps = [
  {
    title: 'Link a spending account',
    description: 'Use Plaid Link to securely connect the card or checking account you spend from most often.',
    icon: CreditCard,
  },
  {
    title: 'Pick a savings target',
    description: 'Tell Kash Out how much you want to auto-save each month with smart rules.',
    icon: Sparkles,
  },
  {
    title: 'Watch the automation',
    description: 'Every new transaction is scanned against your rules so dollars move to the right buckets automatically.',
    icon: CheckCircle2,
  },
]

const Onboarding = () => {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading } = useAuth()
  const [checkingAccounts, setCheckingAccounts] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [notificationStatus, setNotificationStatus] = useState<'idle' | 'loading' | 'success' | 'error' | 'blocked'>(
    'idle'
  )
  const [notificationMessage, setNotificationMessage] = useState<string>('')

  const checkAccounts = useCallback(async () => {
    setCheckingAccounts(true)
    setError(null)

    const [data, err] = await api.financial.getAccounts()

    if (err) {
      setError(err.message || 'Unable to verify your accounts. Try again in a moment.')
      setCheckingAccounts(false)
      return
    }

    const hasAccounts = (data?.accounts?.length ?? 0) > 0
    if (hasAccounts) {
      navigate('/recommendations', { replace: true })
      return
    }

    setCheckingAccounts(false)
  }, [navigate])

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/login', { replace: true })
    }
  }, [isAuthenticated, isLoading, navigate])

  useEffect(() => {
    checkAccounts()
  }, [checkAccounts])

  if (isLoading || checkingAccounts) {
    return (
      <PageShell>
        <PageBody variant="centered">
          <p className="text-sm text-neutral-500">Preparing your onboarding checklist…</p>
        </PageBody>
      </PageShell>
    )
  }

  return (
    <PageShell>
      <PageHeader title="Connect your first account" />

      <PageBody className="space-y-5">
        <Card className="rounded-2xl border border-gray-100 shadow-sm">
          <CardContent className="space-y-4 p-5">
            <div className="space-y-2 text-center">
              <ShieldCheck className="mx-auto h-10 w-10 text-neutral-900" />
              <h1 className="text-2xl font-semibold text-neutral-900">Kick off Kash Out</h1>
              <p className="text-sm text-neutral-500">
                Link at least one bank or card account so we can start watching your transactions in real time.
              </p>
            </div>

            <PlaidLinkButton
              className="w-full rounded-xl"
              onSuccess={checkAccounts}
            >
              Connect a bank via Plaid
            </PlaidLinkButton>

            {error && (
              <div className="rounded-lg border border-red-100 bg-red-50 px-3 py-2 text-sm text-red-600">
                {error}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="rounded-2xl border border-gray-100">
          <CardContent className="space-y-4 p-5">
            {steps.map((step) => (
              <div key={step.title} className="flex gap-4">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-gray-100 text-neutral-900">
                  <step.icon className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-neutral-900">{step.title}</p>
                  <p className="text-sm text-neutral-500">{step.description}</p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="rounded-2xl border border-gray-100">
          <CardContent className="space-y-4 p-5">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gray-100 text-neutral-900">
                <Bell className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-semibold text-neutral-900">Stay notified</p>
                <p className="text-xs text-neutral-500">Get push alerts when rules run or accounts need attention.</p>
              </div>
            </div>
            <div className="space-y-3 text-sm text-neutral-600">
              <p>We’ll only send actionable notifications about rule executions and onboarding reminders.</p>
              <div className="flex flex-col gap-3 sm:flex-row">
                <Button
                  type="button"
                  className="flex-1"
                  disabled={notificationStatus === 'loading' || notificationStatus === 'success'}
                  onClick={async () => {
                    setNotificationStatus('loading')
                    const result = await subscribeToNotifications()
                    setNotificationStatus(result.status)
                    setNotificationMessage(result.message)
                  }}
                >
                  {notificationStatus === 'loading' ? 'Enabling…' : 'Enable notifications'}
                </Button>
                {notificationStatus === 'success' && (
                  <Button type="button" variant="ghost" className="flex-1" disabled>
                    Enabled
                  </Button>
                )}
              </div>
              {notificationStatus !== 'idle' && notificationMessage && (
                <p className={`text-xs ${notificationStatus === 'success' ? 'text-neutral-500' : 'text-red-600'}`}>
                  {notificationMessage}
                </p>
              )}
            </div>
          </CardContent>
        </Card>

        <EmptyState
          title="Need help connecting?"
          description="Reach out to support@kashout.app and we'll walk through your first sync together."
          action={(
            <Button variant="ghost" onClick={checkAccounts}>
              Refresh status
            </Button>
          )}
        />
      </PageBody>
    </PageShell>
  )
}

export default Onboarding
