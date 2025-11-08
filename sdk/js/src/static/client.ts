/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { decode, encode } from 'cbor-x'

import { Konsole } from './logger'
import { DEFAULT_RETRY_OPTIONS, type RetryOptions } from './retry'
import type { Codec, ErrorResponse, Result } from './types'
import { COMPAT, DEFAULT_HEADERS, isNetworkError } from './utils'

/**
 * Context object containing all request-related information
 * Used to pass data between helper methods in a structured way
 */
interface RequestContext<T, R, E extends string> {
  method: string
  path: string
  body?: T
  requestCodec?: Codec<T>
  responseCodec?: Codec<R>
  requestHeaders?: Record<string, string>
  directRequest: boolean
  directResponse: boolean
  isIdempotent: boolean
  attempt: number
  logger: ReturnType<Konsole['groupCollapsed']>
}

/**
 * Abstract base client providing HTTP request functionality with retry logic,
 * comprehensive error handling, and debug logging capabilities
 */
export abstract class BaseClient {
  private readonly fetch: typeof fetch
  private readonly debug: boolean
  private readonly konsole: Konsole = new Konsole()
  private commonHeaders: Record<string, string> = {}
  private readonly retryOptions: RetryOptions

  protected constructor(
    protected readonly endpoint: string,
    options?: {
      headers?: Record<string, string>
      fetch?: typeof fetch
      debug?: boolean
      retryOptions?: Partial<RetryOptions>
    },
  ) {
    this.fetch = options?.fetch ?? globalThis.fetch.bind(globalThis)
    this.debug = options?.debug ?? false
    this.commonHeaders = { ...(options?.headers ?? {}) }
    this.retryOptions = {
      ...DEFAULT_RETRY_OPTIONS,
      ...options?.retryOptions,
    }
  }

  /**
   * Preprocesses ReadableStream body for retry compatibility
   */
  private async preprocessStreamBody<T>(body: T): Promise<Result<T, string>> {
    if (!(body instanceof ReadableStream)) {
      return [body, null]
    }

    if (body.locked) {
      return [null, { code: 'invalid_argument', message: 'stream is locked' }]
    }

    if (!this.retryOptions.retryStreamsAsArrays) {
      return [body, null]
    }

    try {
      const reader = body.getReader()
      const chunks: Uint8Array[] = []
      let totalLength = 0

      // Read all chunks and calculate total length
      let result
      while (!(result = await reader.read()).done) {
        const chunk = result.value
        chunks.push(chunk)
        totalLength += chunk.length
      }

      // Allocate a new Uint8Array with the exact size
      const allBytes = new Uint8Array(totalLength)

      // Copy each chunk into the appropriate position
      let position = 0
      for (const chunk of chunks) {
        allBytes.set(chunk, position)
        position += chunk.length
      }

      return [allBytes as T, null]
    } catch (e) {
      return [null, { code: 'unknown', message: `Failed to read stream: ${(e as Error).message}` }]
    }
  }

  /**
   * Builds the fetch request options including URL, headers, and body
   */
  private buildRequestOptions<T>(context: RequestContext<T, any, any>): { url: string; options: RequestInit } {
    const url = this.endpoint + context.path

    const headers = {
      ...DEFAULT_HEADERS,
      ...(context.directRequest
        ? { 'Content-Type': 'application/octet-stream' }
        : { 'Content-Type': 'application/cbor' }),
      ...this.commonHeaders,
      ...context.requestHeaders,
    } as Record<string, string>

    const options: RequestInit = {
      method: context.method,
      headers,
      credentials: COMPAT.isCredentialsSupported ? 'include' : undefined,
    }

    // Handle request body
    if (context.body !== undefined) {
      if (context.directRequest) {
        options.body = context.body as BodyInit
      } else {
        options.body = encode(context.requestCodec?.encode(context.body)) as BodyInit
      }
    }

    return { url, options }
  }

  /**
   * Debug logging utilities
   */
  private logRequestDetails<T>(
    context: RequestContext<T, any, any>,
    url: string,
    headers: Record<string, string>,
  ): void {
    if (!this.debug) return

    context.logger.time('Duration')
    context.logger.log('%cRequest Details:', 'color: #666')
    context.logger.log('%c→ URL:', 'color: #888', url)
    context.logger.log('%c→ Headers:', 'color: #888', headers)

    if (context.body !== undefined) {
      if (context.directRequest && context.body instanceof ReadableStream) {
        context.logger.log('%c→ Body (Stream):', 'color: #888', 'ReadableStream')
      } else {
        context.logger.log('%c→ Body:', 'color: #888', context.body)
      }
    }
  }

  private logErrorResponse(context: RequestContext<any, any, any>, status: number, error: any): void {
    if (!this.debug) return

    context.logger.log('%cError Response:', 'color: #F44336; font-weight: bold;')
    context.logger.log('%c→ Status:', 'color: #888', status)
    context.logger.log('%c→ Error:', 'color: #888', error)
    context.logger.timeEnd('Duration')
  }

  private logSuccessResponse<R>(
    context: RequestContext<any, R, any>,
    status: number,
    response: any,
    isStream: boolean = false,
  ): void {
    if (!this.debug) return

    const responseType = isStream ? ' (Stream)' : ''
    context.logger.log(`%cSuccess Response${responseType}:`, 'color: #4CAF50; font-weight: bold;')
    context.logger.log('%c→ Status:', 'color: #888', status)

    if (isStream) {
      context.logger.log('%c→ Response:', 'color: #888', 'ReadableStream')
    } else if (response === undefined) {
      context.logger.log('%c→ Response:', 'color: #888', 'No content')
    } else {
      const decodedResponse = context.responseCodec ? context.responseCodec.decode(response) : response
      context.logger.log('%c→ Response:', 'color: #888', decodedResponse)
    }

    context.logger.timeEnd('Duration')
    context.logger.end()
  }

  private logNetworkError(context: RequestContext<any, any, any>, error: any, canRetry: boolean): void {
    if (!this.debug) return

    context.logger.log('%cNetwork Error:', 'color: #FF9800; font-weight: bold;')
    context.logger.log('%c→ Method:', 'color: #888', context.method)
    context.logger.log('%c→ Path:', 'color: #888', context.path)
    context.logger.log('%c→ Attempt:', 'color: #888', `${context.attempt + 1}/${this.retryOptions.maxRetries + 1}`)
    context.logger.log('%c→ Error Type:', 'color: #888', error.name || 'NetworkError')
    context.logger.log('%c→ Error Message:', 'color: #888', error.message)

    if (error.cause) {
      context.logger.log('%c→ Error Cause:', 'color: #888', error.cause)
    }

    if (canRetry) {
      context.logger.log('%c→ Will Retry:', 'color: #888', 'Yes')
      const delay =
        context.attempt === 0
          ? this.retryOptions.initialDelay
          : Math.min(
              this.retryOptions.initialDelay * Math.pow(this.retryOptions.backoffFactor, context.attempt),
              this.retryOptions.maxDelay,
            )
      context.logger.log('%c→ Retry Delay:', 'color: #888', `~${Math.round(delay)}ms`)
    } else {
      context.logger.log('%c→ Will Retry:', 'color: #888', 'No')
      if (!context.isIdempotent) {
        context.logger.log('%c→ Reason:', 'color: #888', 'Non-idempotent method')
      } else if (context.attempt >= this.retryOptions.maxRetries) {
        context.logger.log('%c→ Reason:', 'color: #888', 'Max retries exceeded')
      }
    }

    context.logger.timeEnd('Duration')
    context.logger.end()
  }

  private logUnknownError(context: RequestContext<any, any, any>, error: any): void {
    if (!this.debug) return

    context.logger.log('%cUnknown Error:', 'color: #FF9800; font-weight: bold;')
    context.logger.log('%c→ Method:', 'color: #888', context.method)
    context.logger.log('%c→ Path:', 'color: #888', context.path)
    context.logger.log('%c→ Attempt:', 'color: #888', `${context.attempt + 1}/${this.retryOptions.maxRetries + 1}`)
    context.logger.log('%c→ Error Type:', 'color: #888', error.constructor?.name || 'Error')
    context.logger.log('%c→ Error Message:', 'color: #888', error.message || 'Unknown error occurred')
    context.logger.log('%c→ Error Stack:', 'color: #888', error.stack)

    if (error.cause) {
      context.logger.log('%c→ Error Cause:', 'color: #888', error.cause)
    }

    context.logger.timeEnd('Duration')
    context.logger.end()
  }

  private logStreamError(context: RequestContext<any, any, any>, error: any): void {
    if (!this.debug) return

    context.logger.log('%cNetwork Error with Stream:', 'color: #FF9800; font-weight: bold;')
    context.logger.log('%c→ Method:', 'color: #888', context.method)
    context.logger.log('%c→ Path:', 'color: #888', context.path)
    context.logger.log('%c→ Attempt:', 'color: #888', `${context.attempt + 1}/${this.retryOptions.maxRetries + 1}`)
    context.logger.log('%c→ Error Type:', 'color: #888', error.name || 'NetworkError')
    context.logger.log('%c→ Error Message:', 'color: #888', error.message)
    context.logger.log('%c→ Cannot retry stream after first attempt', 'color: #888')
    context.logger.timeEnd('Duration')
    context.logger.end()
  }

  /**
   * Handles successful HTTP responses
   */
  private async handleSuccessResponse<R, E extends string>(
    context: RequestContext<any, R, E>,
    response: Response,
  ): Promise<Result<R | Response, E>> {
    // Handle direct response - return the Response object directly
    if (context.directResponse) {
      this.logSuccessResponse(context, response.status, null, true)
      return [response as R, null]
    }

    // Handle regular response
    const buffer = await response.arrayBuffer()

    if (buffer.byteLength === 0 && !context.responseCodec) {
      this.logSuccessResponse(context, response.status, undefined)
      return [undefined as R, null]
    }

    const decoded = buffer.byteLength > 0 ? decode(new Uint8Array(buffer)) : undefined
    const result = context.responseCodec ? context.responseCodec.decode(decoded) : (decoded as R)

    this.logSuccessResponse(context, response.status, decoded)
    return [result, null]
  }

  /**
   * Handles HTTP error responses
   */
  private async handleErrorResponse<E extends string>(
    context: RequestContext<any, any, E>,
    response: Response,
  ): Promise<ErrorResponse & { code: E }> {
    const buffer = await response.arrayBuffer()
    const decoded = buffer.byteLength > 0 ? (decode(new Uint8Array(buffer)) as ErrorResponse) : undefined

    this.logErrorResponse(context, response.status, decoded)
    return decoded as ErrorResponse & { code: E }
  }

  /**
   * Determines if a request should be retried based on error type and attempt count
   */
  private shouldRetryRequest<T>(context: RequestContext<T, any, any>, isRetryableStatus: boolean = false): boolean {
    if (!context.isIdempotent || context.attempt >= this.retryOptions.maxRetries) {
      return false
    }

    // Don't retry streaming requests after first attempt
    if (context.body instanceof ReadableStream && context.attempt > 0) {
      return false
    }

    return isRetryableStatus
  }

  /**
   * Calculates and applies retry delay
   */
  private async applyRetryDelay(context: RequestContext<any, any, any>, isRateLimit: boolean = false): Promise<void> {
    if (isRateLimit) {
      await new Promise((resolve) => setTimeout(resolve, this.retryOptions.rateLimitDelay))
    } else {
      // Use exponential backoff for other retryable errors
      const backoff =
        Math.min(
          this.retryOptions.initialDelay * Math.pow(this.retryOptions.backoffFactor, context.attempt),
          this.retryOptions.maxDelay,
        ) *
        (0.8 + Math.random() * 0.4) // Add jitter

      await new Promise((resolve) => setTimeout(resolve, backoff))
    }
  }

  /**
   * Makes an HTTP request with retry logic and comprehensive error handling
   *
   * @param method - HTTP method (GET, POST, etc.)
   * @param path - API endpoint path
   * @param body - Request body data
   * @param requestCodec - Codec for encoding request body
   * @param responseCodec - Codec for decoding response body
   * @param requestHeaders - Additional headers to include
   * @param directRequest - Whether to send body as raw bytes instead of CBOR
   * @param directResponse - Whether to return Response object directly
   * @returns Promise resolving to [data, null] on success or [null, error] on failure
   */
  public async request<T = undefined, R = void, E extends string = string>(
    method: string,
    path: string,
    body?: T,
    requestCodec?: Codec<T>,
    responseCodec?: Codec<R>,
    requestHeaders?: Record<string, string>,
    directRequest?: boolean,
    directResponse?: true,
  ): Promise<Result<Response, E>>

  /**
   * Makes an HTTP request with retry logic and comprehensive error handling
   *
   * @param method - HTTP method (GET, POST, etc.)
   * @param path - API endpoint path
   * @param body - Request body data
   * @param requestCodec - Codec for encoding request body
   * @param responseCodec - Codec for decoding response body
   * @param requestHeaders - Additional headers to include
   * @param directRequest - Whether to send body as raw bytes instead of CBOR
   * @param directResponse - Whether to return Response object directly
   * @returns Promise resolving to [data, null] on success or [null, error] on failure
   */
  public async request<T = undefined, R = void, E extends string = string>(
    method: string,
    path: string,
    body?: T,
    requestCodec?: Codec<T>,
    responseCodec?: Codec<R>,
    requestHeaders?: Record<string, string>,
    directRequest?: boolean,
    directResponse?: false,
  ): Promise<Result<R, E>>

  /**
   * Makes an HTTP request with retry logic and comprehensive error handling
   *
   * @param method - HTTP method (GET, POST, etc.)
   * @param path - API endpoint path
   * @param body - Request body data
   * @param requestCodec - Codec for encoding request body
   * @param responseCodec - Codec for decoding response body
   * @param requestHeaders - Additional headers to include
   * @param directRequest - Whether to send body as raw bytes instead of CBOR
   * @param directResponse - Whether to return Response object directly
   * @returns Promise resolving to [data, null] on success or [null, error] on failure
   */
  public async request<T = undefined, R = void, E extends string = string>(
    method: string,
    path: string,
    body?: T,
    requestCodec?: Codec<T>,
    responseCodec?: Codec<R>,
    requestHeaders?: Record<string, string>,
    directRequest: boolean = false,
    directResponse: boolean = false,
  ): Promise<Result<R | Response, E>> {
    // Preprocess stream body if needed
    if (body instanceof ReadableStream && directRequest) {
      const [processedBody, error] = await this.preprocessStreamBody(body)
      if (error) {
        return [null, { code: 'invalid_argument' as E, message: error.message }]
      }
      body = processedBody
    }

    const isIdempotent = this.retryOptions.idempotentMethods.includes(method.toUpperCase())
    let attempt = 0
    let lastError: Result<R | Response, E> | null = null

    while (attempt <= (isIdempotent ? this.retryOptions.maxRetries : 0)) {
      const logger = this.konsole.groupCollapsed(
        `%c[${new Date().toLocaleTimeString()}] ${method} ${path}${attempt > 0 ? ` (Retry ${attempt}/${this.retryOptions.maxRetries})` : ''}`,
        'color: #2196F3; font-weight: bold;',
      )

      const context: RequestContext<T, R, E> = {
        method,
        path,
        body,
        requestCodec,
        responseCodec,
        requestHeaders,
        directRequest,
        directResponse,
        isIdempotent,
        attempt,
        logger,
      }

      try {
        const { url, options } = this.buildRequestOptions(context)
        this.logRequestDetails(context, url, options.headers as Record<string, string>)

        const response = await this.fetch(url, options)

        // Handle error responses
        if (!response.ok) {
          const errorResponse = await this.handleErrorResponse(context, response)
          const isRetryableStatus = this.retryOptions.retryableStatusCodes.includes(response.status)

          if (this.shouldRetryRequest(context, isRetryableStatus)) {
            lastError = [null, errorResponse]

            if (this.debug) {
              context.logger.log('%c→ Retrying...', 'color: #888')
            }

            context.logger.end()
            await this.applyRetryDelay(context, response.status === 429)
            attempt++
            continue
          }

          context.logger.end()
          return [null, errorResponse]
        }

        // Handle successful response
        return await this.handleSuccessResponse(context, response)
      } catch (error) {
        if (isNetworkError(error)) {
          const networkError = {
            code: 'network' as E,
            message: (error as Error).message,
          }

          // Don't retry streaming requests after first attempt
          if (context.body instanceof ReadableStream && context.attempt > 0) {
            this.logStreamError(context, error)
            context.logger.end()
            return [null, { ...networkError, message: 'Cannot retry stream after network error' }]
          }

          lastError = [null, networkError]
          const canRetry = this.shouldRetryRequest(context)
          this.logNetworkError(context, error, canRetry)

          if (canRetry) {
            context.logger.end()
            await this.applyRetryDelay(context)
            attempt++
            continue
          }

          context.logger.end()
          break
        }

        // Handle unknown errors
        this.logUnknownError(context, error)
        context.logger.end()

        return [
          null,
          {
            code: 'unknown' as E,
            message: (error as Error).message,
          },
        ]
      }
    }

    return lastError!
  }

  /**
   * Set or update a persistent header applied to all subsequent requests
   */
  public setHeader(name: string, value: string): void {
    this.commonHeaders[name] = value
  }

  /**
   * Bulk set/update persistent headers (merged with existing)
   */
  public setHeaders(headers: Record<string, string>): void {
    for (const [k, v] of Object.entries(headers)) {
      this.commonHeaders[k] = v
    }
  }

  /**
   * Remove a persistent header
   */
  public removeHeader(name: string): void {
    delete this.commonHeaders[name]
  }

  /**
   * Clear all persistent headers
   */
  public clearHeaders(): void {
    this.commonHeaders = {}
  }

  /**
   * Get a persistent header value
   */
  public getHeader(name: string): string | undefined {
    return this.commonHeaders[name]
  }

  /**
   * Get a shallow copy of all persistent headers
   */
  public getHeaders(): Record<string, string> {
    return { ...this.commonHeaders }
  }
}
