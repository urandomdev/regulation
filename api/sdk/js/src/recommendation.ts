/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  GetRecommendationsResponse,
  GetRecommendationsResponseCodec,
} from './regulation/server/handlers/recommendation'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class RecommendationAPI {
  constructor(private client: BaseClient) {}

  /**
   * GetRecommendations analyzes user's recent transactions and suggests savings rules
   */
  async getRecommendations(): Promise<Result<GetRecommendationsResponse, 'internal_error' | 'invalid_request'>> {
    return this.client.request<undefined, GetRecommendationsResponse, 'internal_error' | 'invalid_request'>(
      'GET',
      '/recommendations',
      undefined,
      undefined,
      GetRecommendationsResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
