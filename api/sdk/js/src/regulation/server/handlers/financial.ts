/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'
import { Account, AccountCodec } from './models'

export interface GetAccountsResponse {
  accounts: Account[]
}

export const GetAccountsResponseCodec: Codec<GetAccountsResponse> = {
  encode: (data: GetAccountsResponse): object => ({
    accounts: data.accounts ? data.accounts.map((item: any) => AccountCodec.encode(item)) : [],
  }),
  decode: (encoded: any): GetAccountsResponse => ({
    accounts: encoded['accounts'] ? encoded['accounts'].map((item: any) => AccountCodec.decode(item)) : [],
  }),
}
