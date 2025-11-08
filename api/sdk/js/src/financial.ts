/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { GetAccountsResponse, GetAccountsResponseCodec } from './regulation/server/handlers/financial'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class FinancialAPI {
  constructor(private client: BaseClient) {}

  /**
   * GetAccounts retrieves all accounts for the authenticated user
   */
  async getAccounts(): Promise<Result<GetAccountsResponse, 'internal_error'>> {
    return this.client.request<undefined, GetAccountsResponse, 'internal_error'>(
      'GET',
      '/financial/accounts',
      undefined,
      undefined,
      GetAccountsResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
