import { api } from '@/lib/api'

const VAPID_KEY_CACHE_KEY = 'kashout:vapid'

function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(base64)
  const outputArray = new Uint8Array(rawData.length)
  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i)
  }
  return outputArray
}

export interface NotificationSetupResult {
  status: 'success' | 'error' | 'blocked'
  message: string
}

export async function subscribeToNotifications(): Promise<NotificationSetupResult> {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
    return { status: 'error', message: 'Notifications are not supported in this browser.' }
  }

  if (!('Notification' in window)) {
    return { status: 'error', message: 'Notifications API unavailable.' }
  }

  const permission = await Notification.requestPermission()
  if (permission !== 'granted') {
    return { status: 'blocked', message: 'Enable notifications in your browser to stay updated.' }
  }

  await navigator.serviceWorker.register('/sw.js')
  const readyRegistration = await navigator.serviceWorker.ready

  let subscription = await readyRegistration.pushManager.getSubscription()
  if (!subscription) {
    const vapidKey = await getVapidKey()
    if (!vapidKey) {
      return { status: 'error', message: 'Unable to fetch notification key.' }
    }

    subscription = await readyRegistration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(vapidKey),
    })
  }

  const subscriptionJson = subscription.toJSON()
  const endpoint = subscriptionJson.endpoint
  const keys = subscriptionJson.keys || {}
  if (!endpoint || !keys.p256dh || !keys.auth) {
    return { status: 'error', message: 'Incomplete subscription details.' }
  }

  const [, err] = await api.notification.subscribe({
    endpoint,
    p256dh: keys.p256dh,
    auth: keys.auth,
  })

  if (err) {
    return { status: 'error', message: err.message || 'Failed to register notifications.' }
  }

  return { status: 'success', message: 'Notifications enabled. We will alert you when rules fire.' }
}

async function getVapidKey(): Promise<string | null> {
  const cached = sessionStorage.getItem(VAPID_KEY_CACHE_KEY)
  if (cached) {
    return cached
  }

  const [data, err] = await api.notification.getVAPIDPublicKey()
  if (err || !data?.public_key) {
    return null
  }

  sessionStorage.setItem(VAPID_KEY_CACHE_KEY, data.public_key)
  return data.public_key
}
