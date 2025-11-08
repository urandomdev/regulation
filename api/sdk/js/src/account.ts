/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { SignupRequest, SignupRequestCodec } from './regulation/server/handlers/account'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class AccountAPI {
  constructor(private client: BaseClient) {}

  /**
   * Signup creates a new user account
   */
  async signup(params: SignupRequest): Promise<Result<void, 'internal_error'>> {
    return this.client.request<SignupRequest, void, 'internal_error'>(
      'POST',
      '/account/signup',
      params,
      SignupRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
