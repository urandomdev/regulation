/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import {
  CreateRuleRequest,
  CreateRuleRequestCodec,
  CreateRuleResponse,
  CreateRuleResponseCodec,
  ListRulesResponse,
  ListRulesResponseCodec,
  RuleExecutionsResponse,
  RuleExecutionsResponseCodec,
  RuleResponse,
  RuleResponseCodec,
  UpdateRuleRequest,
  UpdateRuleRequestCodec,
} from './regulation/server/handlers/rule'
import { BaseClient } from './static/client'
import { type Codec, type Result } from './static/types'

export class RuleAPI {
  constructor(private client: BaseClient) {}

  /**
   * CreateRule creates a new savings rule
   */
  async createRule(
    params: CreateRuleRequest,
  ): Promise<Result<CreateRuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<
      CreateRuleRequest,
      CreateRuleResponse,
      'internal_error' | 'invalid_parameters' | 'not_found'
    >(
      'POST',
      '/rules',
      params,
      CreateRuleRequestCodec,
      CreateRuleResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * DeleteRule deletes a rule
   */
  async deleteRule(id: string): Promise<Result<void, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<undefined, void, 'internal_error' | 'invalid_parameters' | 'not_found'>(
      'DELETE',
      `/rules/${id}`,
      undefined,
      undefined,
      undefined,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * GetRule retrieves a specific rule by ID
   */
  async getRule(id: string): Promise<Result<RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<undefined, RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>(
      'GET',
      `/rules/${id}`,
      undefined,
      undefined,
      RuleResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * GetRuleExecutions retrieves execution history for a specific rule
   */
  async getRuleExecutions(
    id: string,
  ): Promise<Result<RuleExecutionsResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<
      undefined,
      RuleExecutionsResponse,
      'internal_error' | 'invalid_parameters' | 'not_found'
    >(
      'GET',
      `/rules/${id}/executions`,
      undefined,
      undefined,
      RuleExecutionsResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * ListRules retrieves all rules for the current user
   */
  async listRules(): Promise<Result<ListRulesResponse, 'internal_error'>> {
    return this.client.request<undefined, ListRulesResponse, 'internal_error'>(
      'GET',
      '/rules',
      undefined,
      undefined,
      ListRulesResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * ToggleRule toggles a rule's active status
   */
  async toggleRule(id: string): Promise<Result<RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<undefined, RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>(
      'PATCH',
      `/rules/${id}/toggle`,
      undefined,
      undefined,
      RuleResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }

  /**
   * UpdateRule updates an existing rule
   */
  async updateRule(
    id: string,
    params: UpdateRuleRequest,
  ): Promise<Result<RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>> {
    return this.client.request<UpdateRuleRequest, RuleResponse, 'internal_error' | 'invalid_parameters' | 'not_found'>(
      'PATCH',
      `/rules/${id}`,
      params,
      UpdateRuleRequestCodec,
      RuleResponseCodec,
      undefined, // headers
      false, // direct request
      false, // direct response
    )
  }
}
