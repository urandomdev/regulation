/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  LoginRequest,
  LoginRequestCodec,
  SignupRequest,
  SignupRequestCodec,
  UserResponse,
  UserResponseCodec,
} from './regulation/server/handlers/account'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class AccountAPI {
  constructor(private client: BaseClient) {}

  /**
   * Login authenticates a user and creates a session
   */
  async login(params: LoginRequest): Promise<Result<void, 'internal_error'>> {
    return this.client.request<LoginRequest, void, 'internal_error'>(
      'POST',
      '/account/login',
      params,
      LoginRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * Logout destroys the current session
   */
  async logout(): Promise<Result<void, 'internal_error'>> {
    return this.client.request<undefined, void, 'internal_error'>(
      'POST',
      '/account/logout',
      undefined,
      undefined,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * Me returns the currently authenticated user's information
   */
  async me(): Promise<Result<UserResponse, 'internal_error'>> {
    return this.client.request<undefined, UserResponse, 'internal_error'>(
      'GET',
      '/account/me',
      undefined,
      undefined,
      UserResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

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
