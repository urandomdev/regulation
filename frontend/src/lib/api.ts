import { API } from '@urandomdev/regulation'

// Get API URL from environment variable or default to localhost
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Create and export the API client instance
export const api = new API(API_URL, {
  debug: import.meta.env.DEV, // Enable debug mode in development
})

// Export types for convenience
export type { Result } from '@urandomdev/regulation/static/types'
