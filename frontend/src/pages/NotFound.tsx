import { Link } from 'react-router'
import { ArrowLeft } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { PageShell, PageBody } from '@/components/layout/Page'

const NotFound = () => {
  return (
    <PageShell>
      <PageBody variant="centered" className="w-full max-w-md text-center">
        <Card className="w-full border border-gray-100 shadow-sm">
          <CardContent className="space-y-4 p-6">
            <p className="text-sm font-semibold tracking-wide text-amber-600">404 error</p>
            <h1 className="text-2xl font-bold text-neutral-900">Page not found</h1>
            <p className="text-sm text-neutral-500">
              The page you're looking for might have been moved or no longer exists.
            </p>
            <div className="flex flex-col gap-3 sm:flex-row">
              <Button asChild className="w-full">
                <Link to="/">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Back to home
                </Link>
              </Button>
              <Button asChild variant="ghost" className="w-full">
                <Link to="/dashboard">Open dashboard</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </PageBody>
    </PageShell>
  )
}

export default NotFound
