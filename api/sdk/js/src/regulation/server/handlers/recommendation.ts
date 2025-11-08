/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { Codec } from '../../../static/types'

export interface GetRecommendationsResponse {
  suggestions: RuleSuggestionResponse[]
  overall_analysis: string
  priority_suggestion: string
  analysis_period_days: bigint
  transaction_count: bigint
  total_spent: number
}
export interface RuleSuggestionResponse {
  name: string
  category: string
  action_type: string
  action_value: number
  min_amount_cents?: bigint | null
  max_amount_cents?: bigint | null
  estimated_savings: number
  confidence: string
  reasoning: string
  impact_level: string
}

export const GetRecommendationsResponseCodec: Codec<GetRecommendationsResponse> = {
  encode: (data: GetRecommendationsResponse): object => ({
    suggestions: data.suggestions ? data.suggestions.map((item: any) => RuleSuggestionResponseCodec.encode(item)) : [],
    overall_analysis: data.overall_analysis,
    priority_suggestion: data.priority_suggestion,
    analysis_period_days: data.analysis_period_days,
    transaction_count: data.transaction_count,
    total_spent: data.total_spent,
  }),
  decode: (encoded: any): GetRecommendationsResponse => ({
    suggestions: encoded['suggestions']
      ? encoded['suggestions'].map((item: any) => RuleSuggestionResponseCodec.decode(item))
      : [],
    overall_analysis: encoded['overall_analysis'],
    priority_suggestion: encoded['priority_suggestion'],
    analysis_period_days: BigInt(encoded['analysis_period_days']),
    transaction_count: BigInt(encoded['transaction_count']),
    total_spent: encoded['total_spent'],
  }),
}

export const RuleSuggestionResponseCodec: Codec<RuleSuggestionResponse> = {
  encode: (data: RuleSuggestionResponse): object => ({
    name: data.name,
    category: data.category,
    action_type: data.action_type,
    action_value: data.action_value,
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
    estimated_savings: data.estimated_savings,
    confidence: data.confidence,
    reasoning: data.reasoning,
    impact_level: data.impact_level,
  }),
  decode: (encoded: any): RuleSuggestionResponse => ({
    name: encoded['name'],
    category: encoded['category'],
    action_type: encoded['action_type'],
    action_value: encoded['action_value'],
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
    estimated_savings: encoded['estimated_savings'],
    confidence: encoded['confidence'],
    reasoning: encoded['reasoning'],
    impact_level: encoded['impact_level'],
  }),
}
