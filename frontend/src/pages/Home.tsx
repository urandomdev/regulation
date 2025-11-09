import { Link, Navigate } from 'react-router'
import { Sparkles, ShieldCheck, LineChart, ArrowRight } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { PageShell, PageBody } from '@/components/layout/Page'
import { useAuth } from '@/hooks/useAuth'

const features: Array<{ title: string; description: string; icon: LucideIcon }> = [
  {
    title: 'Automated rules',
    description: 'Round up charges, split spend, and route cash to the right accounts without spreadsheets.',
    icon: Sparkles,
  },
  {
    title: 'Bank-level guardrails',
    description: 'Use Plaid Link to connect accounts securely. Kash Out never stores your credentials.',
    icon: ShieldCheck,
  },
  {
    title: 'Cashflow clarity',
    description: 'Track income, spend, and savings progress each month in a single mobile-friendly view.',
    icon: LineChart,
  },
]

const Home = () => {
  const { isAuthenticated, isLoading } = useAuth()

  if (!isLoading && isAuthenticated) {
    return <Navigate to="/dashboard" replace />
  }

  return (
    <PageShell>
      <PageBody className="max-w-2xl space-y-10 text-center">
        <div className="space-y-4">
          <p className="text-sm font-semibold uppercase tracking-wide text-teal-600">Who Are You To Spend?</p>
          <h1 className="text-4xl font-bold text-neutral-900">Kash Out</h1>
          <p className="text-lg text-neutral-600">
            A rule-based personal finance companion that keeps your spending intentional and your savings on autopilot.
          </p>
        </div>

        <Card className="rounded-3xl border border-gray-100 shadow-sm text-left">
          <CardContent className="space-y-6 p-6">
            {features.map((feature) => (
              <div key={feature.title} className="flex gap-4">
                <div className="h-12 w-12 rounded-2xl bg-gray-100 text-teal-600 flex items-center justify-center">
                  <feature.icon className="h-5 w-5" />
                </div>
                <div>
                  <h2 className="text-base font-semibold text-neutral-900">{feature.title}</h2>
                  <p className="text-sm text-neutral-500">{feature.description}</p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <div className="flex flex-col gap-3 sm:flex-row">
          <Button asChild className="w-full">
            <Link to="/signup">Create account</Link>
          </Button>
          <Button asChild variant="secondary" className="w-full">
            <Link to="/login">Sign in</Link>
          </Button>
        </div>

        <Button asChild variant="ghost" className="mx-auto text-sm text-neutral-500 hover:text-neutral-900">
          <Link to="/dashboard" className="flex items-center gap-1">
            Peek at the dashboard
            <ArrowRight className="h-4 w-4" />
          </Link>
        </Button>
      </PageBody>
    </PageShell>
  )
}

export default Home
