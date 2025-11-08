/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { UUID } from '@deltalaboratory/uuid'

import { Codec } from '../../../static/types'
import { ActionType, ActionTypeCodec, Category, CategoryCodec } from '../../internal/ent/rule'

export interface CreateRuleRequest {
  name: string
  category: Category
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  action_type: ActionType
  action_value: number
  target_account_id: UUID
  priority: bigint
}
export interface CreateRuleResponse {
  id: UUID
  name: string
  category: string
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  action_type: string
  action_value: number
  target_account_id: UUID
  priority: bigint
  is_active: boolean
  execution_count: bigint
  total_saved_cents: bigint
  created_at: any
  updated_at: any
}
export interface ListRulesResponse {
  rules: RuleResponse[]
  total: bigint
}
export interface RuleExecutionResponse {
  id: UUID
  rule_id: UUID
  transaction_id: UUID
  amount_cents: bigint
  source_account_id: UUID
  target_account_id: UUID
  status: string
  error_message?: string | null
  created_at: any
  completed_at?: any
  transaction?: TransactionSummary | null
}
export interface RuleExecutionsResponse {
  executions: RuleExecutionResponse[]
  total: bigint
}
export interface RuleResponse {
  id: UUID
  name: string
  category: string
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  action_type: string
  action_value: number
  target_account_id: UUID
  priority: bigint
  is_active: boolean
  execution_count: bigint
  total_saved_cents: bigint
  created_at: any
  updated_at: any
}
export interface TransactionSummary {
  name: string
  merchant_name?: string | null
  amount: bigint
  date: any
  category: string
}
export interface UpdateRuleRequest {
  name?: string | null
  category?: Category | null
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  action_type?: ActionType | null
  action_value?: number | null
  priority?: bigint | null
}

export const CreateRuleRequestCodec: Codec<CreateRuleRequest> = {
  encode: (data: CreateRuleRequest): object => ({
    name: data.name,
    category: CategoryCodec.encode(data.category),
    ...(data.min_amount_cents !== undefined && data.min_amount_cents !== null
      ? {
          min_amount_cents:
            data.min_amount_cents !== undefined && data.min_amount_cents !== null ? data.min_amount_cents : null,
        }
      : {}),
    ...(data.max_amount_cents !== undefined && data.max_amount_cents !== null
      ? {
          max_amount_cents:
            data.max_amount_cents !== undefined && data.max_amount_cents !== null ? data.max_amount_cents : null,
        }
      : {}),
    action_type: ActionTypeCodec.encode(data.action_type),
    action_value: data.action_value,
    target_account_id: data.target_account_id.buffer,
    priority: data.priority,
  }),
  decode: (encoded: any): CreateRuleRequest => ({
    name: encoded['name'],
    category: CategoryCodec.decode(encoded['category']),
    min_amount_cents:
      encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
        ? encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
          ? BigInt(encoded['min_amount_cents'])
          : null
        : undefined,
    max_amount_cents:
      encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
        ? encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
          ? BigInt(encoded['max_amount_cents'])
          : null
        : undefined,
    action_type: ActionTypeCodec.decode(encoded['action_type']),
    action_value: encoded['action_value'],
    target_account_id: new UUID(encoded['target_account_id']),
    priority: BigInt(encoded['priority']),
  }),
}

export const CreateRuleResponseCodec: Codec<CreateRuleResponse> = {
  encode: (data: CreateRuleResponse): object => ({
    id: data.id.buffer,
    name: data.name,
    category: data.category,
    ...(data.min_amount_cents !== undefined && data.min_amount_cents !== null
      ? {
          min_amount_cents:
            data.min_amount_cents !== undefined && data.min_amount_cents !== null ? data.min_amount_cents : null,
        }
      : {}),
    ...(data.max_amount_cents !== undefined && data.max_amount_cents !== null
      ? {
          max_amount_cents:
            data.max_amount_cents !== undefined && data.max_amount_cents !== null ? data.max_amount_cents : null,
        }
      : {}),
    action_type: data.action_type,
    action_value: data.action_value,
    target_account_id: data.target_account_id.buffer,
    priority: data.priority,
    is_active: data.is_active,
    execution_count: data.execution_count,
    total_saved_cents: data.total_saved_cents,
    created_at: data.created_at,
    updated_at: data.updated_at,
  }),
  decode: (encoded: any): CreateRuleResponse => ({
    id: new UUID(encoded['id']),
    name: encoded['name'],
    category: encoded['category'],
    min_amount_cents:
      encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
        ? encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
          ? BigInt(encoded['min_amount_cents'])
          : null
        : undefined,
    max_amount_cents:
      encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
        ? encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
          ? BigInt(encoded['max_amount_cents'])
          : null
        : undefined,
    action_type: encoded['action_type'],
    action_value: encoded['action_value'],
    target_account_id: new UUID(encoded['target_account_id']),
    priority: BigInt(encoded['priority']),
    is_active: encoded['is_active'],
    execution_count: BigInt(encoded['execution_count']),
    total_saved_cents: BigInt(encoded['total_saved_cents']),
    created_at: encoded['created_at'],
    updated_at: encoded['updated_at'],
  }),
}

export const ListRulesResponseCodec: Codec<ListRulesResponse> = {
  encode: (data: ListRulesResponse): object => ({
    rules: data.rules ? data.rules.map((item: any) => RuleResponseCodec.encode(item)) : [],
    total: data.total,
  }),
  decode: (encoded: any): ListRulesResponse => ({
    rules: encoded['rules'] ? encoded['rules'].map((item: any) => RuleResponseCodec.decode(item)) : [],
    total: BigInt(encoded['total']),
  }),
}

export const RuleExecutionResponseCodec: Codec<RuleExecutionResponse> = {
  encode: (data: RuleExecutionResponse): object => ({
    id: data.id.buffer,
    rule_id: data.rule_id.buffer,
    transaction_id: data.transaction_id.buffer,
    amount_cents: data.amount_cents,
    source_account_id: data.source_account_id.buffer,
    target_account_id: data.target_account_id.buffer,
    status: data.status,
    ...(data.error_message !== undefined && data.error_message !== null ? { error_message: data.error_message } : {}),
    created_at: data.created_at,
    ...(data.completed_at !== undefined && data.completed_at !== null ? { completed_at: data.completed_at } : {}),
    ...(data.transaction !== undefined && data.transaction !== null
      ? { transaction: TransactionSummaryCodec.encode(data.transaction) }
      : {}),
  }),
  decode: (encoded: any): RuleExecutionResponse => ({
    id: new UUID(encoded['id']),
    rule_id: new UUID(encoded['rule_id']),
    transaction_id: new UUID(encoded['transaction_id']),
    amount_cents: BigInt(encoded['amount_cents']),
    source_account_id: new UUID(encoded['source_account_id']),
    target_account_id: new UUID(encoded['target_account_id']),
    status: encoded['status'],
    error_message:
      encoded['error_message'] !== undefined && encoded['error_message'] !== null
        ? encoded['error_message']
        : undefined,
    created_at: encoded['created_at'],
    completed_at: encoded['completed_at'] !== undefined ? encoded['completed_at'] : undefined,
    transaction:
      encoded['transaction'] !== undefined && encoded['transaction'] !== null
        ? TransactionSummaryCodec.decode(encoded['transaction'])
        : undefined,
  }),
}

export const RuleExecutionsResponseCodec: Codec<RuleExecutionsResponse> = {
  encode: (data: RuleExecutionsResponse): object => ({
    executions: data.executions ? data.executions.map((item: any) => RuleExecutionResponseCodec.encode(item)) : [],
    total: data.total,
  }),
  decode: (encoded: any): RuleExecutionsResponse => ({
    executions: encoded['executions']
      ? encoded['executions'].map((item: any) => RuleExecutionResponseCodec.decode(item))
      : [],
    total: BigInt(encoded['total']),
  }),
}

export const RuleResponseCodec: Codec<RuleResponse> = {
  encode: (data: RuleResponse): object => ({
    id: data.id.buffer,
    name: data.name,
    category: data.category,
    ...(data.min_amount_cents !== undefined && data.min_amount_cents !== null
      ? {
          min_amount_cents:
            data.min_amount_cents !== undefined && data.min_amount_cents !== null ? data.min_amount_cents : null,
        }
      : {}),
    ...(data.max_amount_cents !== undefined && data.max_amount_cents !== null
      ? {
          max_amount_cents:
            data.max_amount_cents !== undefined && data.max_amount_cents !== null ? data.max_amount_cents : null,
        }
      : {}),
    action_type: data.action_type,
    action_value: data.action_value,
    target_account_id: data.target_account_id.buffer,
    priority: data.priority,
    is_active: data.is_active,
    execution_count: data.execution_count,
    total_saved_cents: data.total_saved_cents,
    created_at: data.created_at,
    updated_at: data.updated_at,
  }),
  decode: (encoded: any): RuleResponse => ({
    id: new UUID(encoded['id']),
    name: encoded['name'],
    category: encoded['category'],
    min_amount_cents:
      encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
        ? encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
          ? BigInt(encoded['min_amount_cents'])
          : null
        : undefined,
    max_amount_cents:
      encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
        ? encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
          ? BigInt(encoded['max_amount_cents'])
          : null
        : undefined,
    action_type: encoded['action_type'],
    action_value: encoded['action_value'],
    target_account_id: new UUID(encoded['target_account_id']),
    priority: BigInt(encoded['priority']),
    is_active: encoded['is_active'],
    execution_count: BigInt(encoded['execution_count']),
    total_saved_cents: BigInt(encoded['total_saved_cents']),
    created_at: encoded['created_at'],
    updated_at: encoded['updated_at'],
  }),
}

export const TransactionSummaryCodec: Codec<TransactionSummary> = {
  encode: (data: TransactionSummary): object => ({
    name: data.name,
    ...(data.merchant_name !== undefined && data.merchant_name !== null ? { merchant_name: data.merchant_name } : {}),
    amount: data.amount,
    date: data.date,
    category: data.category,
  }),
  decode: (encoded: any): TransactionSummary => ({
    name: encoded['name'],
    merchant_name:
      encoded['merchant_name'] !== undefined && encoded['merchant_name'] !== null
        ? encoded['merchant_name']
        : undefined,
    amount: BigInt(encoded['amount']),
    date: encoded['date'],
    category: encoded['category'],
  }),
}

export const UpdateRuleRequestCodec: Codec<UpdateRuleRequest> = {
  encode: (data: UpdateRuleRequest): object => ({
    ...(data.name !== undefined && data.name !== null ? { name: data.name } : {}),
    ...(data.category !== undefined && data.category !== null ? { category: CategoryCodec.encode(data.category) } : {}),
    ...(data.min_amount_cents !== undefined && data.min_amount_cents !== null
      ? {
          min_amount_cents:
            data.min_amount_cents !== undefined && data.min_amount_cents !== null ? data.min_amount_cents : null,
        }
      : {}),
    ...(data.max_amount_cents !== undefined && data.max_amount_cents !== null
      ? {
          max_amount_cents:
            data.max_amount_cents !== undefined && data.max_amount_cents !== null ? data.max_amount_cents : null,
        }
      : {}),
    ...(data.action_type !== undefined && data.action_type !== null
      ? { action_type: ActionTypeCodec.encode(data.action_type) }
      : {}),
    ...(data.action_value !== undefined && data.action_value !== null ? { action_value: data.action_value } : {}),
    ...(data.priority !== undefined && data.priority !== null
      ? { priority: data.priority !== undefined && data.priority !== null ? data.priority : null }
      : {}),
  }),
  decode: (encoded: any): UpdateRuleRequest => ({
    name: encoded['name'] !== undefined && encoded['name'] !== null ? encoded['name'] : undefined,
    category:
      encoded['category'] !== undefined && encoded['category'] !== null
        ? CategoryCodec.decode(encoded['category'])
        : undefined,
    min_amount_cents:
      encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
        ? encoded['min_amount_cents'] !== undefined && encoded['min_amount_cents'] !== null
          ? BigInt(encoded['min_amount_cents'])
          : null
        : undefined,
    max_amount_cents:
      encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
        ? encoded['max_amount_cents'] !== undefined && encoded['max_amount_cents'] !== null
          ? BigInt(encoded['max_amount_cents'])
          : null
        : undefined,
    action_type:
      encoded['action_type'] !== undefined && encoded['action_type'] !== null
        ? ActionTypeCodec.decode(encoded['action_type'])
        : undefined,
    action_value:
      encoded['action_value'] !== undefined && encoded['action_value'] !== null ? encoded['action_value'] : undefined,
    priority:
      encoded['priority'] !== undefined && encoded['priority'] !== null
        ? encoded['priority'] !== undefined && encoded['priority'] !== null
          ? BigInt(encoded['priority'])
          : null
        : undefined,
  }),
}
