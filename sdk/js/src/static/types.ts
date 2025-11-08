/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.

export type Result<T, E extends string = string> =
  | [T, null]
  | [
      null,
      ErrorResponse & {
        code: 'invalid_argument' | 'unknown' | 'network' | E
      },
    ]

export interface ErrorResponse {
  code: string
  message: string
  meta?: Record<string, unknown>
}

export interface Codec<T, E = object> {
  encode: (data: T) => E
  decode: (encoded: any) => T
}
