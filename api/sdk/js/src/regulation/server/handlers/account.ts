/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'

export interface SignupRequest {
  email: string
  password: string
  nickname: string
}

export const SignupRequestCodec: Codec<SignupRequest> = {
  encode: (data: SignupRequest): object => ({
    email: data.email,
    password: data.password,
    nickname: data.nickname,
  }),
  decode: (encoded: any): SignupRequest => ({
    email: encoded['email'],
    password: encoded['password'],
    nickname: encoded['nickname'],
  }),
}
