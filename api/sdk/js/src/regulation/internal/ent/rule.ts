/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'

export type ActionType = 'multiply' | 'fixed'
export type Category =
  | 'Dining'
  | 'Groceries'
  | 'Transport'
  | 'Shopping'
  | 'Subscriptions'
  | 'Entertainment'
  | 'Bills'
  | 'Misc'

export const ActionTypeCodec: Codec<ActionType, string> = {
  encode: (data: ActionType): string => data,
  decode: (encoded: any): ActionType => encoded as ActionType,
}

export const CategoryCodec: Codec<Category, string> = {
  encode: (data: Category): string => data,
  decode: (encoded: any): Category => encoded as Category,
}
