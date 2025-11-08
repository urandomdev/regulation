/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  GetAccountTransactionsRequest,
  GetAccountTransactionsRequestCodec,
  GetAccountTransactionsResponse,
  GetAccountTransactionsResponseCodec,
  GetAccountsResponse,
  GetAccountsResponseCodec,
  GetCashflowRequest,
  GetCashflowRequestCodec,
  GetCashflowResponse,
  GetCashflowResponseCodec,
  GetTransactionsRequest,
  GetTransactionsRequestCodec,
  GetTransactionsResponse,
  GetTransactionsResponseCodec,
} from './regulation/server/handlers/financial'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class FinancialAPI {
  constructor(private client: BaseClient) {}

  /**
   * GetAccountTransactions retrieves transactions for a specific account
   */
  async getAccountTransactions(
    id: string,
    params: GetAccountTransactionsRequest,
  ): Promise<Result<GetAccountTransactionsResponse, 'internal_error'>> {
    return this.client.request<GetAccountTransactionsRequest, GetAccountTransactionsResponse, 'internal_error'>(
      'POST',
      `/financial/accounts/${id}/transactions`,
      params,
      GetAccountTransactionsRequestCodec,
      GetAccountTransactionsResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

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

  /**
   * GetCashflow retrieves income, spend, and NET for a date range
   */
  async getCashflow(params: GetCashflowRequest): Promise<Result<GetCashflowResponse, 'internal_error'>> {
    return this.client.request<GetCashflowRequest, GetCashflowResponse, 'internal_error'>(
      'POST',
      '/financial/cashflow',
      params,
      GetCashflowRequestCodec,
      GetCashflowResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * GetTransactions retrieves transactions for the authenticated user with filtering and pagination
   */
  async getTransactions(params: GetTransactionsRequest): Promise<Result<GetTransactionsResponse, 'internal_error'>> {
    return this.client.request<GetTransactionsRequest, GetTransactionsResponse, 'internal_error'>(
      'POST',
      '/financial/transactions',
      params,
      GetTransactionsRequestCodec,
      GetTransactionsResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
