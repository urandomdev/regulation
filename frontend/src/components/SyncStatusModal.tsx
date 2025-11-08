import { Button } from '@/components/ui/button'

type SyncPhase = 'idle' | 'syncing' | 'success' | 'error'

interface LatestActivity {
  id: string
  name: string
  amountCents: number
  date: string
}

interface SyncStatusModalProps {
  open: boolean
  phase: SyncPhase
  message: string
  onClose: () => void
  activities?: LatestActivity[]
  totalSynced?: number
}

const phaseTitles: Record<SyncPhase, string> = {
  idle: 'Preparing sync',
  syncing: 'Sync in progress',
  success: 'Latest activity captured',
  error: 'Sync incomplete',
}

function formatAmount(amountCents: number) {
  const absolute = Math.abs(amountCents) / 100
  const sign = amountCents < 0 ? '+' : '-'
  return `${sign}$${absolute.toFixed(2)}`
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}

export function SyncStatusModal({ open, phase, message, onClose, activities, totalSynced }: SyncStatusModalProps) {
  if (!open) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/40" aria-hidden="true" />
      <div className="relative z-10 w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
        <div className="space-y-4">
          <div>
            <p className="text-sm font-semibold text-teal-600">{phaseTitles[phase]}</p>
            <p className="mt-1 text-sm text-neutral-600">{message}</p>
            {typeof totalSynced === 'number' && totalSynced > 0 && (
              <p className="mt-1 text-xs text-neutral-500">
                Synced {totalSynced} transaction{totalSynced === 1 ? '' : 's'} so far
              </p>
            )}
          </div>

          {(phase === 'syncing' || phase === 'idle') && (
            <div className="flex items-center gap-3 text-sm text-neutral-500">
              <span className="h-4 w-4 animate-spin rounded-full border-2 border-gray-200 border-t-teal-500" />
              Syncing with your institution…
            </div>
          )}

          {activities?.length ? (
            <div className="rounded-xl border border-gray-100 bg-gray-50 p-4">
              <p className="text-xs uppercase tracking-wide text-neutral-500">Latest synced activity</p>
              <div className="mt-3 space-y-3">
                {activities.map((activity) => (
                  <div key={activity.id} className="flex items-center justify-between gap-4">
                    <div>
                      <p className="text-sm font-semibold text-neutral-900">{activity.name}</p>
                      <p className="text-xs text-neutral-500">{formatDate(activity.date)}</p>
                    </div>
                    <p className={`text-sm font-semibold ${activity.amountCents < 0 ? 'text-emerald-600' : 'text-neutral-900'}`}>
                      {formatAmount(activity.amountCents)}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          ) : null}

          {phase === 'error' && (
            <div className="rounded-xl border border-red-100 bg-red-50 px-3 py-2 text-sm text-red-600">
              We connected your account but couldn’t pull transactions this time. You can retry from Accounts.
            </div>
          )}

          <div className="flex justify-end">
            <Button
              variant={phase === 'syncing' ? 'secondary' : 'default'}
              disabled={phase === 'syncing'}
              onClick={onClose}
            >
              {phase === 'syncing' ? 'Syncing…' : 'Dismiss'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}

export type { LatestActivity, SyncPhase }
