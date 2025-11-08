/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { UUID } from '@deltalaboratory/uuid'

import { Codec } from '../../../static/types'
import { Account, AccountCodec, Transaction, TransactionCodec } from './models'

export interface GetAccountTransactionsRequest {
  start?: Date | null
  end?: Date | null
  limit: bigint
  offset: bigint
}
export interface GetAccountTransactionsResponse {
  transactions: Transaction[]
  total: bigint
}
export interface GetAccountsResponse {
  accounts: Account[]
}
export interface GetCashflowRequest {
  start?: Date | null
  end?: Date | null
}
export interface GetCashflowResponse {
  total_income: bigint
  total_spend: bigint
  net: bigint
  start: Date
  end: Date
}
export interface GetTransactionsRequest {
  start?: Date | null
  end?: Date | null
  account_id?: UUID | null
  limit: bigint
  offset: bigint
}
export interface GetTransactionsResponse {
  transactions: Transaction[]
  total: bigint
}

export const GetAccountTransactionsRequestCodec: Codec<GetAccountTransactionsRequest> = {
  encode: (data: GetAccountTransactionsRequest): object => ({
    ...(data.start !== undefined && data.start !== null
      ? { start: data.start !== undefined && data.start !== null ? data.start.getTime() / 1000 : null }
      : {}),
    ...(data.end !== undefined && data.end !== null
      ? { end: data.end !== undefined && data.end !== null ? data.end.getTime() / 1000 : null }
      : {}),
    limit: data.limit,
    offset: data.offset,
  }),
  decode: (encoded: any): GetAccountTransactionsRequest => ({
    start:
      encoded['start'] !== undefined && encoded['start'] !== null
        ? encoded['start'] !== undefined && encoded['start'] !== null
          ? new Date(encoded['start'] * 1000)
          : null
        : undefined,
    end:
      encoded['end'] !== undefined && encoded['end'] !== null
        ? encoded['end'] !== undefined && encoded['end'] !== null
          ? new Date(encoded['end'] * 1000)
          : null
        : undefined,
    limit: BigInt(encoded['limit']),
    offset: BigInt(encoded['offset']),
  }),
}

export const GetAccountTransactionsResponseCodec: Codec<GetAccountTransactionsResponse> = {
  encode: (data: GetAccountTransactionsResponse): object => ({
    transactions: data.transactions ? data.transactions.map((item: any) => TransactionCodec.encode(item)) : [],
    total: data.total,
  }),
  decode: (encoded: any): GetAccountTransactionsResponse => ({
    transactions: encoded['transactions']
      ? encoded['transactions'].map((item: any) => TransactionCodec.decode(item))
      : [],
    total: BigInt(encoded['total']),
  }),
}

export const GetAccountsResponseCodec: Codec<GetAccountsResponse> = {
  encode: (data: GetAccountsResponse): object => ({
    accounts: data.accounts ? data.accounts.map((item: any) => AccountCodec.encode(item)) : [],
  }),
  decode: (encoded: any): GetAccountsResponse => ({
    accounts: encoded['accounts'] ? encoded['accounts'].map((item: any) => AccountCodec.decode(item)) : [],
  }),
}

export const GetCashflowRequestCodec: Codec<GetCashflowRequest> = {
  encode: (data: GetCashflowRequest): object => ({
    ...(data.start !== undefined && data.start !== null
      ? { start: data.start !== undefined && data.start !== null ? data.start.getTime() / 1000 : null }
      : {}),
    ...(data.end !== undefined && data.end !== null
      ? { end: data.end !== undefined && data.end !== null ? data.end.getTime() / 1000 : null }
      : {}),
  }),
  decode: (encoded: any): GetCashflowRequest => ({
    start:
      encoded['start'] !== undefined && encoded['start'] !== null
        ? encoded['start'] !== undefined && encoded['start'] !== null
          ? new Date(encoded['start'] * 1000)
          : null
        : undefined,
    end:
      encoded['end'] !== undefined && encoded['end'] !== null
        ? encoded['end'] !== undefined && encoded['end'] !== null
          ? new Date(encoded['end'] * 1000)
          : null
        : undefined,
  }),
}

export const GetCashflowResponseCodec: Codec<GetCashflowResponse> = {
  encode: (data: GetCashflowResponse): object => ({
    total_income: data.total_income,
    total_spend: data.total_spend,
    net: data.net,
    start: data.start.getTime() / 1000,
    end: data.end.getTime() / 1000,
  }),
  decode: (encoded: any): GetCashflowResponse => ({
    total_income: BigInt(encoded['total_income']),
    total_spend: BigInt(encoded['total_spend']),
    net: BigInt(encoded['net']),
    start: new Date(encoded['start'] * 1000),
    end: new Date(encoded['end'] * 1000),
  }),
}

export const GetTransactionsRequestCodec: Codec<GetTransactionsRequest> = {
  encode: (data: GetTransactionsRequest): object => ({
    ...(data.start !== undefined && data.start !== null
      ? { start: data.start !== undefined && data.start !== null ? data.start.getTime() / 1000 : null }
      : {}),
    ...(data.end !== undefined && data.end !== null
      ? { end: data.end !== undefined && data.end !== null ? data.end.getTime() / 1000 : null }
      : {}),
    ...(data.account_id !== undefined && data.account_id !== null
      ? { account_id: data.account_id !== undefined && data.account_id !== null ? data.account_id.buffer : null }
      : {}),
    limit: data.limit,
    offset: data.offset,
  }),
  decode: (encoded: any): GetTransactionsRequest => ({
    start:
      encoded['start'] !== undefined && encoded['start'] !== null
        ? encoded['start'] !== undefined && encoded['start'] !== null
          ? new Date(encoded['start'] * 1000)
          : null
        : undefined,
    end:
      encoded['end'] !== undefined && encoded['end'] !== null
        ? encoded['end'] !== undefined && encoded['end'] !== null
          ? new Date(encoded['end'] * 1000)
          : null
        : undefined,
    account_id:
      encoded['account_id'] !== undefined && encoded['account_id'] !== null
        ? encoded['account_id'] !== undefined && encoded['account_id'] !== null
          ? new UUID(encoded['account_id'])
          : null
        : undefined,
    limit: BigInt(encoded['limit']),
    offset: BigInt(encoded['offset']),
  }),
}

export const GetTransactionsResponseCodec: Codec<GetTransactionsResponse> = {
  encode: (data: GetTransactionsResponse): object => ({
    transactions: data.transactions ? data.transactions.map((item: any) => TransactionCodec.encode(item)) : [],
    total: data.total,
  }),
  decode: (encoded: any): GetTransactionsResponse => ({
    transactions: encoded['transactions']
      ? encoded['transactions'].map((item: any) => TransactionCodec.decode(item))
      : [],
    total: BigInt(encoded['total']),
  }),
}
