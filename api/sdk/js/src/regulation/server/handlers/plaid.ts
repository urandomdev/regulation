/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { UUID } from '@deltalaboratory/uuid'

import { Codec } from '../../../static/types'

export interface CreateLinkTokenResponse {
  link_token: string
}
export interface DisconnectAccountResponse {
  success: boolean
}
export interface ExchangeTokenRequest {
  public_token: string
}
export interface SyncTransactionsRequest {
  item_id?: UUID | null
}

export const CreateLinkTokenResponseCodec: Codec<CreateLinkTokenResponse> = {
  encode: (data: CreateLinkTokenResponse): object => ({
    link_token: data.link_token,
  }),
  decode: (encoded: any): CreateLinkTokenResponse => ({
    link_token: encoded['link_token'],
  }),
}

export const DisconnectAccountResponseCodec: Codec<DisconnectAccountResponse> = {
  encode: (data: DisconnectAccountResponse): object => ({
    success: data.success,
  }),
  decode: (encoded: any): DisconnectAccountResponse => ({
    success: encoded['success'],
  }),
}

export const ExchangeTokenRequestCodec: Codec<ExchangeTokenRequest> = {
  encode: (data: ExchangeTokenRequest): object => ({
    public_token: data.public_token,
  }),
  decode: (encoded: any): ExchangeTokenRequest => ({
    public_token: encoded['public_token'],
  }),
}

export const SyncTransactionsRequestCodec: Codec<SyncTransactionsRequest> = {
  encode: (data: SyncTransactionsRequest): object => ({
    ...(data.item_id !== undefined && data.item_id !== null
      ? { item_id: data.item_id !== undefined && data.item_id !== null ? data.item_id.buffer : null }
      : {}),
  }),
  decode: (encoded: any): SyncTransactionsRequest => ({
    item_id:
      encoded['item_id'] !== undefined && encoded['item_id'] !== null
        ? encoded['item_id'] !== undefined && encoded['item_id'] !== null
          ? new UUID(encoded['item_id'])
          : null
        : undefined,
  }),
}
