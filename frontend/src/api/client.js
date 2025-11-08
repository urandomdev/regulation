import { API } from '@urandomdev/regulation';

// Initialize API client
const apiClient = new API(import.meta.env.VITE_API_URL || 'http://localhost:8080', {
  headers: {},
  debug: import.meta.env.DEV
});

// Export APIs
export const financialApi = apiClient.financial;
export const accountApi = apiClient.account;

// Utility functions
export const utils = {
  /**
   * Convert cents to dollars
   */
  centsToDollars: (cents) => {
    return cents / 100;
  },

  /**
   * Format date as YYYY-MM-DD
   */
  formatDate: (date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  },

  /**
   * Get current month date range (first day to last day)
   */
  getCurrentMonthRange: () => {
    const now = new Date();
    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1);
    const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    return {
      start: utils.formatDate(firstDay),
      end: utils.formatDate(lastDay),
    };
  },
};
