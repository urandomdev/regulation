/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.

export const DEFAULT_HEADERS = {} as const

export const COMPAT = {
  // Cloudflare Workers: The 'credentials' field on 'RequestInitializerDict' is not implemented
  // https://github.com/cloudflare/workers-sdk/issues/2514
  isCredentialsSupported: 'credentials' in Request.prototype,
} as const

export const isNetworkError = (error: unknown): boolean => {
  const isError = (value: unknown): boolean => Object.prototype.toString.call(value) === '[object Error]'

  const NETWORK_ERROR_MESSAGES: string[] = [
    'network error', // Chrome
    'Failed to fetch', // Chrome
    'NetworkError when attempting to fetch resource.', // Firefox
    'The Internet connection appears to be offline.', // Safari 16
    'Load failed', // Safari 17+
    'Network request failed', // `cross-fetch`
    'fetch failed', // Undici (Node.js)
    'terminated',
  ]

  if (
    !error ||
    !isError(error) ||
    (error as Error).name !== 'TypeError' ||
    typeof (error as Error).message !== 'string'
  ) {
    return false
  }

  // We do an extra check for Safari 17+ as it has a very generic error message.
  // Network errors in Safari have no stack.
  if ((error as Error).message === 'Load failed') {
    return (error as Error).stack === undefined
  }

  return NETWORK_ERROR_MESSAGES.includes((error as Error).message)
}
