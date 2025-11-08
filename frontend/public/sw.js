self.addEventListener('push', event => {
  if (!event.data) {
    return
  }

  const payload = event.data.json()
  const title = payload.title || 'Regulation update'
  const options = {
    body: payload.body || 'Open the app to see more details.',
    icon: payload.icon || '/vite.svg',
    data: payload.url || '/',
  }

  event.waitUntil(self.registration.showNotification(title, options))
})

self.addEventListener('notificationclick', event => {
  event.notification.close()
  const targetUrl = event.notification.data || '/'
  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then(clientList => {
      for (const client of clientList) {
        if ('focus' in client) {
          client.focus()
          client.navigate(targetUrl)
          return
        }
      }
      if (clients.openWindow) {
        return clients.openWindow(targetUrl)
      }
    })
  )
})
