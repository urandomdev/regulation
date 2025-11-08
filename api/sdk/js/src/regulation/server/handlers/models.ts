/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { UUID } from '@deltalaboratory/uuid'

import { Codec } from '../../../static/types'

export interface Account {
  id: UUID
  name: string
  type: string
  mask?: string | null
  current_balance: bigint
  available_balance?: bigint | null
  is_active: boolean
}
export interface Transaction {
  id: UUID
  account_id: UUID
  amount: bigint
  date: Date
  name: string
  merchant_name?: string | null
  category: string
  pending: boolean
  payment_channel?: string | null
}

export const AccountCodec: Codec<Account> = {
  encode: (data: Account): object => ({
    id: data.id.buffer,
    name: data.name,
    type: data.type,
    ...(data.mask !== undefined && data.mask !== null ? { mask: data.mask } : {}),
    current_balance: data.current_balance,
    ...(data.available_balance !== undefined && data.available_balance !== null
      ? {
          available_balance:
            data.available_balance !== undefined && data.available_balance !== null ? data.available_balance : null,
        }
      : {}),
    is_active: data.is_active,
  }),
  decode: (encoded: any): Account => ({
    id: new UUID(encoded['id']),
    name: encoded['name'],
    type: encoded['type'],
    mask: encoded['mask'] !== undefined && encoded['mask'] !== null ? encoded['mask'] : undefined,
    current_balance: BigInt(encoded['current_balance']),
    available_balance:
      encoded['available_balance'] !== undefined && encoded['available_balance'] !== null
        ? encoded['available_balance'] !== undefined && encoded['available_balance'] !== null
          ? BigInt(encoded['available_balance'])
          : null
        : undefined,
    is_active: encoded['is_active'],
  }),
}

export const TransactionCodec: Codec<Transaction> = {
  encode: (data: Transaction): object => ({
    id: data.id.buffer,
    account_id: data.account_id.buffer,
    amount: data.amount,
    date: data.date.getTime() / 1000,
    name: data.name,
    ...(data.merchant_name !== undefined && data.merchant_name !== null ? { merchant_name: data.merchant_name } : {}),
    category: data.category,
    pending: data.pending,
    ...(data.payment_channel !== undefined && data.payment_channel !== null
      ? { payment_channel: data.payment_channel }
      : {}),
  }),
  decode: (encoded: any): Transaction => ({
    id: new UUID(encoded['id']),
    account_id: new UUID(encoded['account_id']),
    amount: BigInt(encoded['amount']),
    date: new Date(encoded['date'] * 1000),
    name: encoded['name'],
    merchant_name:
      encoded['merchant_name'] !== undefined && encoded['merchant_name'] !== null
        ? encoded['merchant_name']
        : undefined,
    category: encoded['category'],
    pending: encoded['pending'],
    payment_channel:
      encoded['payment_channel'] !== undefined && encoded['payment_channel'] !== null
        ? encoded['payment_channel']
        : undefined,
  }),
}
