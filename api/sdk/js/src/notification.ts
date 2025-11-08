/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  SubscribeRequest,
  SubscribeRequestCodec,
  UnsubscribeRequest,
  UnsubscribeRequestCodec,
  VAPIDPublicKeyResponse,
  VAPIDPublicKeyResponseCodec,
} from './regulation/server/handlers/notification'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class NotificationAPI {
  constructor(private client: BaseClient) {}

  /**
   * GetVAPIDPublicKey returns the VAPID public key for web push notifications
   */
  async getVAPIDPublicKey(): Promise<Result<VAPIDPublicKeyResponse, 'internal_error' | 'invalid_request'>> {
    return this.client.request<undefined, VAPIDPublicKeyResponse, 'internal_error' | 'invalid_request'>(
      'GET',
      '/notification/vapid',
      undefined,
      undefined,
      VAPIDPublicKeyResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * Subscribe registers a new web push subscription for the authenticated user
   */
  async subscribe(params: SubscribeRequest): Promise<Result<void, 'internal_error' | 'invalid_request'>> {
    return this.client.request<SubscribeRequest, void, 'internal_error' | 'invalid_request'>(
      'POST',
      '/notification/subscribe',
      params,
      SubscribeRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * Unsubscribe removes a web push subscription for the authenticated user
   */
  async unsubscribe(params: UnsubscribeRequest): Promise<Result<void, 'internal_error' | 'invalid_request'>> {
    return this.client.request<UnsubscribeRequest, void, 'internal_error' | 'invalid_request'>(
      'DELETE',
      '/notification/subscribe',
      params,
      UnsubscribeRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
