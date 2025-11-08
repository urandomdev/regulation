/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { AccountAPI } from './account'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class API extends BaseClient {
  account: AccountAPI

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
  }
}
