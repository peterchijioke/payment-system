import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

// Server-side only API client
// This file should only be imported in Server Components, API Routes, or Server Actions
// Using environment variables ensures the backend URL is never exposed to the client

const API_BASE_URL = process.env.API_URL;

if (!API_BASE_URL) {
  throw new Error(
    'API_URL environment variable is not defined. ' +
    'Please add API_URL to your .env.local file.'
  );
}

/**
 * Creates a configured axios instance for server-side API calls
 * The base URL is kept secure on the server and never exposed to clients
 */
function createServerApiClient(): AxiosInstance {
  const client = axios.create({
    baseURL: `${API_BASE_URL}/api/v1`,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor for adding auth tokens or other headers
  client.interceptors.request.use(
    (config) => {
      // Add any server-side only headers here
      // For example: internal API keys or auth tokens
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor for error handling
  client.interceptors.response.use(
    (response) => response,
    (error) => {
      // Server-side error logging
      if (error.response) {
        console.error(`API Error: ${error.response.status} - ${error.response.statusText}`);
      } else if (error.request) {
        console.error('API Error: No response received from server');
      } else {
        console.error('API Error:', error.message);
      }
      return Promise.reject(error);
    }
  );

  return client;
}

// Singleton instance for reuse
let apiClient: AxiosInstance | null = null;

/**
 * Get the server-side API client instance
 * Creates a new instance if one doesn't exist
 */
export function getServerApiClient(): AxiosInstance {
  if (!apiClient) {
    apiClient = createServerApiClient();
  }
  return apiClient;
}

/**
 * Payment API methods - all server-side secure
 */
export const paymentApi = {
  /**
   * Get all payments with pagination and filtering
   */
  getPayments: async (params?: {
    limit?: number;
    offset?: number;
    status?: string;
    start_date?: string;
    end_date?: string;
  }) => {
    const client = getServerApiClient();
    const response = await client.get('/payments', { params });
    return response.data;
  },

  /**
   * Get a single payment by ID
   */
  getPayment: async (id: string) => {
    const client = getServerApiClient();
    const response = await client.get(`/payments/${id}`);
    return response.data;
  },

  /**
   * Create a new payment
   */

  /**
   * Create a new payment
   */
  createPayment: async (data: {
    account_id: string;
    amount: number;
    currency: string;
    destination_currency: string;
    recipient_name: string;
    recipient_account: string;
    recipient_bank: string;
    recipient_country: string;
    reference?: string;
    idempotency_key: string;
  }) => {
    const client = getServerApiClient();
    const response = await client.post('/payments', data, {
      headers: {
        'Idempotency-Key': data.idempotency_key,
      },
    });
    return response.data;
  },
};

/**
 * Type for the API response from /payments endpoint
 */
export interface PaymentsResponse {
  data: {
    limit: number;
    offset: number;
    payments: Payment[];
    total: number;
  };
  success: boolean;
}

export interface Payment {
  transaction: {
    id: string;
    transaction_reference: string;
    idempotency_key: string;
    account_id: string;
    counterparty_id: string | null;
    type: string;
    status: string;
    amount: number;
    currency: string;
    settled_amount: number | null;
    fx_quote_id: string;
    fx_rate: number;
    fx_amount: number;
    fx_currency: string;
    description: string;
    reference: string;
    metadata: any;
    initiated_at: string;
    processed_at: string | null;
    settled_at: string | null;
    completed_at: string | null;
    failed_at: string | null;
    failure_reason: string;
    reversal_reason: string;
    reversed_by_id: string | null;
    version: number;
    created_at: string;
    updated_at: string;
  };
  ledger_entries: any[];
  timeline: { status: string; timestamp: string }[];
}

export default paymentApi;
