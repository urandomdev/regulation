import { financialApi, utils } from '../api/client';
import { UUID } from '@deltalaboratory/uuid';

// Financial data service using real backend API
class FinancialService {
  constructor() {
    this.cache = {
      accounts: null,
      transactions: null,
      cashflow: null,
    };
  }

  // Get all user accounts
  async getAccounts() {
    try {
      const [response, error] = await financialApi.getAccounts();
      if (error) {
        throw new Error(error.message || 'Failed to fetch accounts');
      }
      this.cache.accounts = response.accounts.map(account => ({
        id: account.id,
        name: account.name,
        balance: utils.centsToDollars(Number(account.current_balance)),
        availableBalance: account.available_balance
          ? utils.centsToDollars(Number(account.available_balance))
          : null,
        type: account.type,
        mask: account.mask,
        isActive: account.is_active,
      }));
      return this.cache.accounts;
    } catch (error) {
      console.error('Failed to fetch accounts:', error);
      throw error;
    }
  }

  // Get transactions with optional filters
  async getTransactions(params = {}) {
    try {
      const [response, error] = await financialApi.getTransactions({
        start: params.start ? new Date(params.start) : null,
        end: params.end ? new Date(params.end) : null,
        account_id: params.accountId ? new UUID(params.accountId) : null,
        limit: BigInt(params.limit || 50),
        offset: BigInt(params.offset || 0),
      });

      if (error) {
        throw new Error(error.message || 'Failed to fetch transactions');
      }

      const transactions = response.transactions.map(t => ({
        id: t.id,
        accountId: t.account_id,
        date: t.date,
        name: t.name,
        merchantName: t.merchant_name,
        category: t.category,
        amount: Math.abs(utils.centsToDollars(Number(t.amount))),
        isIncome: t.amount < 0, // Negative amount = income
        pending: t.pending,
        paymentChannel: t.payment_channel,
      }));

      this.cache.transactions = {
        data: transactions,
        total: Number(response.total),
      };

      return this.cache.transactions;
    } catch (error) {
      console.error('Failed to fetch transactions:', error);
      throw error;
    }
  }

  // Get transactions for specific account
  async getAccountTransactions(accountId, params = {}) {
    try {
      const [response, error] = await financialApi.getAccountTransactions(accountId, {
        start: params.start ? new Date(params.start) : null,
        end: params.end ? new Date(params.end) : null,
        limit: BigInt(params.limit || 50),
        offset: BigInt(params.offset || 0),
      });

      if (error) {
        throw new Error(error.message || 'Failed to fetch account transactions');
      }

      return {
        data: response.transactions.map(t => ({
          id: t.id,
          accountId: t.account_id,
          date: t.date,
          name: t.name,
          merchantName: t.merchant_name,
          category: t.category,
          amount: Math.abs(utils.centsToDollars(Number(t.amount))),
          isIncome: t.amount < 0,
          pending: t.pending,
          paymentChannel: t.payment_channel,
        })),
        total: Number(response.total),
      };
    } catch (error) {
      console.error('Failed to fetch account transactions:', error);
      throw error;
    }
  }

  // Get cashflow data for date range
  async getCashflow(start, end) {
    try {
      const [response, error] = await financialApi.getCashflow({
        start: start ? new Date(start) : null,
        end: end ? new Date(end) : null
      });

      if (error) {
        throw new Error(error.message || 'Failed to fetch cashflow');
      }

      this.cache.cashflow = {
        income: utils.centsToDollars(Number(response.total_income)),
        spend: utils.centsToDollars(Number(response.total_spend)),
        net: utils.centsToDollars(Number(response.net)),
        start: response.start,
        end: response.end,
      };

      return this.cache.cashflow;
    } catch (error) {
      console.error('Failed to fetch cashflow:', error);
      throw error;
    }
  }

  // Get current month cashflow
  async getCurrentMonthCashflow() {
    const { start, end } = utils.getCurrentMonthRange();
    return this.getCashflow(start, end);
  }

  // Get recent transactions (last 10 by default)
  async getRecentTransactions(limit = 10) {
    const result = await this.getTransactions({ limit, offset: 0 });
    return result.data;
  }

  // Clear cache
  clearCache() {
    this.cache = {
      accounts: null,
      transactions: null,
      cashflow: null,
    };
  }

  // Get total balance across all accounts
  async getTotalBalance() {
    const accounts = await this.getAccounts();
    return accounts.reduce((sum, account) => sum + account.balance, 0);
  }

  // Get monthly income
  async getMonthlyIncome() {
    const cashflow = await this.getCurrentMonthCashflow();
    return cashflow.income;
  }

  // Get monthly spending
  async getMonthlySpend() {
    const cashflow = await this.getCurrentMonthCashflow();
    return cashflow.spend;
  }

  // Get net cashflow
  async getMonthlyNet() {
    const cashflow = await this.getCurrentMonthCashflow();
    return cashflow.net;
  }
}

// Export singleton instance
export const financialService = new FinancialService();
