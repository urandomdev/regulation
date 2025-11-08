/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { UUID } from '@deltalaboratory/uuid'

import { Codec } from '../../../static/types'

export interface LoginRequest {
  email: string
  password: string
}
export interface SignupRequest {
  email: string
  password: string
  nickname: string
}
export interface UserResponse {
  id: UUID
  email: string
  nickname: string
}

export const LoginRequestCodec: Codec<LoginRequest> = {
  encode: (data: LoginRequest): object => ({
    email: data.email,
    password: data.password,
  }),
  decode: (encoded: any): LoginRequest => ({
    email: encoded['email'],
    password: encoded['password'],
  }),
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

export const UserResponseCodec: Codec<UserResponse> = {
  encode: (data: UserResponse): object => ({
    id: data.id.buffer,
    email: data.email,
    nickname: data.nickname,
  }),
  decode: (encoded: any): UserResponse => ({
    id: new UUID(encoded['id']),
    email: encoded['email'],
    nickname: encoded['nickname'],
  }),
}
