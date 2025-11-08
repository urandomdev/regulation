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
