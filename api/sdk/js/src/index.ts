/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { AccountAPI } from './account'
import { FinancialAPI } from './financial'
import { PlaidAPI } from './plaid'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class API extends BaseClient {
  account: AccountAPI
  financial: FinancialAPI
  plaid: PlaidAPI

  constructor(
    endpoint: string,
    options?: {
      headers?: Record<string, string>
      fetch?: typeof fetch
      debug?: boolean
    },
  ) {
    super(endpoint, options)
    this.account = new AccountAPI(this)
    this.financial = new FinancialAPI(this)
    this.plaid = new PlaidAPI(this)
  }
}
