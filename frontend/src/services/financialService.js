import { financialApi, utils } from '../api/client';

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
        balance: utils.centsToDollars(account.currentBalance),
        availableBalance: account.availableBalance
          ? utils.centsToDollars(account.availableBalance)
          : null,
        type: account.type,
        mask: account.mask,
        isActive: account.isActive,
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
        start: params.start,
        end: params.end,
        accountId: params.accountId,
        limit: params.limit || 50,
        offset: params.offset || 0,
      });

      if (error) {
        throw new Error(error.message || 'Failed to fetch transactions');
      }

      const transactions = response.transactions.map(t => ({
        id: t.id,
        accountId: t.accountId,
        date: t.date,
        name: t.name,
        merchantName: t.merchantName,
        category: t.category,
        amount: Math.abs(utils.centsToDollars(t.amount)),
        isIncome: t.amount < 0, // Negative amount = income
        pending: t.pending,
        paymentChannel: t.paymentChannel,
      }));

      this.cache.transactions = {
        data: transactions,
        total: response.total,
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
        start: params.start,
        end: params.end,
        limit: params.limit || 50,
        offset: params.offset || 0,
      });

      if (error) {
        throw new Error(error.message || 'Failed to fetch account transactions');
      }

      return {
        data: response.transactions.map(t => ({
          id: t.id,
          accountId: t.accountId,
          date: t.date,
          name: t.name,
          merchantName: t.merchantName,
          category: t.category,
          amount: Math.abs(utils.centsToDollars(t.amount)),
          isIncome: t.amount < 0,
          pending: t.pending,
          paymentChannel: t.paymentChannel,
        })),
        total: response.total,
      };
    } catch (error) {
      console.error('Failed to fetch account transactions:', error);
      throw error;
    }
  }

  // Get cashflow data for date range
  async getCashflow(start, end) {
    try {
      const [response, error] = await financialApi.getCashflow({ start, end });

      if (error) {
        throw new Error(error.message || 'Failed to fetch cashflow');
      }

      this.cache.cashflow = {
        income: utils.centsToDollars(response.total_income),
        spend: utils.centsToDollars(response.total_spend),
        net: utils.centsToDollars(response.net),
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
