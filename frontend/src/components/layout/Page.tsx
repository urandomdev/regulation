import { type ReactNode } from 'react'
import { useNavigate } from 'react-router'
import { ChevronLeft } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type BackTarget = string | (() => void)

interface PageShellProps {
  children: ReactNode
  background?: 'muted' | 'plain'
  className?: string
}

export function PageShell({ children, background = 'muted', className }: PageShellProps) {
  return (
    <div
      className={cn(
        'min-h-screen',
        background === 'muted'
          ? 'bg-gray-50 text-gray-900'
          : 'bg-background text-foreground',
        className
      )}
    >
      {children}
    </div>
  )
}

interface PageHeaderProps {
  title?: string
  backTo?: BackTarget
  action?: ReactNode
  sticky?: boolean
}

export function PageHeader({ title, backTo, action, sticky = true }: PageHeaderProps) {
  const navigate = useNavigate()

  const showBackButton = typeof backTo !== 'undefined'

  const handleBack = () => {
    if (typeof backTo === 'function') {
      backTo()
      return
    }

    if (typeof backTo === 'string') {
      navigate(backTo)
      return
    }

    navigate(-1)
  }

  if (!title && !showBackButton && !action) {
    return null
  }

  return (
    <header
      className={cn(
        'flex items-center justify-between border-b border-gray-100 bg-white px-5 py-3',
        sticky && 'sticky top-0 z-10'
      )}
    >
      <div className="flex items-center gap-3">
        {showBackButton && (
          <Button
            variant="ghost"
            size="icon"
            className="text-neutral-400 hover:text-neutral-600 transition-colors"
            onClick={handleBack}
          >
            <ChevronLeft className="h-5 w-5" />
          </Button>
        )}
        {title && <div className="text-base font-semibold text-neutral-800">{title}</div>}
      </div>
      {action}
    </header>
  )
}

interface PageBodyProps {
  children: ReactNode
  variant?: 'stacked' | 'centered'
  className?: string
}

export function PageBody({ children, variant = 'stacked', className }: PageBodyProps) {
  return (
    <main
      className={cn(
        'mx-auto w-full max-w-md p-4',
        variant === 'stacked' && 'space-y-4',
        variant === 'centered' &&
          'flex min-h-[calc(100vh-7rem)] flex-col items-center justify-center gap-4',
        className
      )}
    >
      {children}
    </main>
  )
}

interface EmptyStateProps {
  title: string
  description?: string
  action?: ReactNode
}

export function EmptyState({ title, description, action }: EmptyStateProps) {
  return (
    <div className="text-center">
      <h3 className="text-base font-semibold text-neutral-900">{title}</h3>
      {description && <p className="mt-1 text-sm text-neutral-500">{description}</p>}
      {action && <div className="mt-4 flex justify-center">{action}</div>}
    </div>
  )
}
