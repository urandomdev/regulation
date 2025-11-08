import { API } from '@generated/js';

// Configuration
const API_ENDPOINT = import.meta.env.VITE_API_ENDPOINT || 'http://localhost:8080';
const DEBUG_MODE = import.meta.env.DEV;

// Create and export the API client instance
export const apiClient = new API(API_ENDPOINT, {
  debug: DEBUG_MODE,
  headers: {
    'X-Client-Version': '1.0.0',
  },
});

// Export API for convenience
export { API };