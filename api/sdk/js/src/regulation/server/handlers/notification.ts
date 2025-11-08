/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'

export interface SubscribeRequest {
  endpoint: string
  p256dh: string
  auth: string
}
export interface UnsubscribeRequest {
  endpoint: string
}
export interface VAPIDPublicKeyResponse {
  public_key: string
}

export const SubscribeRequestCodec: Codec<SubscribeRequest> = {
  encode: (data: SubscribeRequest): object => ({
    endpoint: data.endpoint,
    p256dh: data.p256dh,
    auth: data.auth,
  }),
  decode: (encoded: any): SubscribeRequest => ({
    endpoint: encoded['endpoint'],
    p256dh: encoded['p256dh'],
    auth: encoded['auth'],
  }),
}

export const UnsubscribeRequestCodec: Codec<UnsubscribeRequest> = {
  encode: (data: UnsubscribeRequest): object => ({
    endpoint: data.endpoint,
  }),
  decode: (encoded: any): UnsubscribeRequest => ({
    endpoint: encoded['endpoint'],
  }),
}

export const VAPIDPublicKeyResponseCodec: Codec<VAPIDPublicKeyResponse> = {
  encode: (data: VAPIDPublicKeyResponse): object => ({
    public_key: data.public_key,
  }),
  decode: (encoded: any): VAPIDPublicKeyResponse => ({
    public_key: encoded['public_key'],
  }),
}
