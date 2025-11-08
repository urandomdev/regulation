/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'

export interface SubscribeRequest {
  endpoint: string
  p256dh: string
  auth: string
}

export const SubscribeRequestCodec: Codec<SubscribeRequest> = {
  encode: (data: SubscribeRequest): object => ({
    endpoint: data.endpoint,
    p256dh: data.p256dh,
    auth: data.auth,
  }),
  decode: (data: any): SubscribeRequest => ({
    endpoint: data.endpoint,
    p256dh: data.p256dh,
    auth: data.auth,
  }),
}

export interface UnsubscribeRequest {
  endpoint: string
}

export const UnsubscribeRequestCodec: Codec<UnsubscribeRequest> = {
  encode: (data: UnsubscribeRequest): object => ({
    endpoint: data.endpoint,
  }),
  decode: (data: any): UnsubscribeRequest => ({
    endpoint: data.endpoint,
  }),
}

export interface VAPIDPublicKeyResponse {
  publicKey: string
}

export const VAPIDPublicKeyResponseCodec: Codec<VAPIDPublicKeyResponse> = {
  encode: (data: VAPIDPublicKeyResponse): object => ({
    public_key: data.publicKey,
  }),
  decode: (data: any): VAPIDPublicKeyResponse => ({
    publicKey: data.public_key,
  }),
}
