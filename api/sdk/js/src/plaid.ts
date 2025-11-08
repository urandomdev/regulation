/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  CreateLinkTokenResponse,
  CreateLinkTokenResponseCodec,
  DisconnectAccountResponse,
  DisconnectAccountResponseCodec,
  ExchangeTokenRequest,
  ExchangeTokenRequestCodec,
  SyncTransactionsRequest,
  SyncTransactionsRequestCodec,
} from './regulation/server/handlers/plaid'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class PlaidAPI {
  constructor(private client: BaseClient) {}

  /**
   * CreateLinkToken creates a Plaid Link token for account linking
   */
  async createLinkToken(): Promise<Result<CreateLinkTokenResponse, 'internal_error'>> {
    return this.client.request<undefined, CreateLinkTokenResponse, 'internal_error'>(
      'POST',
      '/plaid/create-link-token',
      undefined,
      undefined,
      CreateLinkTokenResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * DisconnectAccount soft-deletes an account (used during onboarding to remove incorrectly linked accounts)
   */
  async disconnectAccount(
    id: string,
  ): Promise<Result<DisconnectAccountResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<
      undefined,
      DisconnectAccountResponse,
      'internal_error' | 'invalid_parameters' | 'not_found'
    >(
      'DELETE',
      `/plaid/accounts/${id}`,
      undefined,
      undefined,
      DisconnectAccountResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * ExchangeToken exchanges a public token for an access token and creates Item and Account records
   */
  async exchangeToken(params: ExchangeTokenRequest): Promise<Result<void, 'internal_error'>> {
    return this.client.request<ExchangeTokenRequest, void, 'internal_error'>(
      'POST',
      '/plaid/exchange-token',
      params,
      ExchangeTokenRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * SyncTransactions manually triggers transaction synchronization
   */
  async syncTransactions(params: SyncTransactionsRequest): Promise<Result<void, 'internal_error'>> {
    return this.client.request<SyncTransactionsRequest, void, 'internal_error'>(
      'POST',
      '/plaid/sync-transactions',
      params,
      SyncTransactionsRequestCodec,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
