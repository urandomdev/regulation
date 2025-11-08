/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.

export interface RetryOptions {
  /**
   * The maximum number of times to retry a request.
   * @default 3
   */
  maxRetries: number
  /**
   * The initial delay (in milliseconds) before the first retry.
   * @default 100
   */
  initialDelay: number
  /**
   * The maximum delay (in milliseconds) between retries.
   * @default 5000
   */
  maxDelay: number
  /**
   * The factor by which the delay increases after each retry.
   * @default 2
   */
  backoffFactor: number
  /**
   * An array of HTTP methods that are considered safe to retry (e.g., GET, HEAD, PUT, DELETE, OPTIONS, TRACE).
   * @default ["GET", "HEAD", "PUT", "DELETE", "OPTIONS", "TRACE"]
   */
  idempotentMethods: string[]
  /**
   * Whether to retry requests with `ReadableStream` bodies by buffering the stream into a `Uint8Array`.
   * This is necessary because streams can only be read once.  Enabling this will consume the stream and retry
   * with the buffered data.
   * @default false
   */
  retryStreamsAsArrays: boolean
  /**
   * HTTP status codes that should trigger a retry.
   * @default [408, 429, 500, 502, 503, 504]
   */
  retryableStatusCodes: number[]
  /**
   * The delay (in milliseconds) to wait after receiving a 429 (rate limit) response.
   * @default 10000
   */
  rateLimitDelay: number
}

export const DEFAULT_RETRY_OPTIONS: RetryOptions = {
  maxRetries: 3,
  initialDelay: 100,
  maxDelay: 5000,
  backoffFactor: 2,
  idempotentMethods: ['GET', 'HEAD', 'PUT', 'DELETE', 'OPTIONS', 'TRACE', 'QUERY'],
  retryStreamsAsArrays: false,
  retryableStatusCodes: [408, 429, 500, 502, 503, 504],
  rateLimitDelay: 10000,
} as const
