import { Account, Transaction, TransactionEvent, LedgerEntry, TransactionStatus, Currency, FXQuote } from './types';

// Mock accounts
export const mockAccounts: Account[] = [
  { id: 'acc-001', name: 'Operating Account USD', accountNumber: 'US1234567890', currency: 'USD', balance: 1250000.00 },
  { id: 'acc-002', name: 'EUR Treasury', accountNumber: 'EU9876543210', currency: 'EUR', balance: 850000.00 },
  { id: 'acc-003', name: 'GBP Operations', accountNumber: 'GB1122334455', currency: 'GBP', balance: 420000.00 },
  { id: 'acc-004', name: 'NGN Local', accountNumber: 'NGN9988776655', currency: 'NGN', balance: 150000000.00 },
];

// Mock transactions
export const mockTransactions: Transaction[] = [
  {
    id: 'txn-001',
    reference: 'TXN-2024-001234',
    sender: mockAccounts[0],
    recipient: { name: 'Acme Corp Ltd', accountNumber: 'GB5544332211', country: 'United Kingdom', currency: 'GBP' },
    amount: 50000,
    sourceCurrency: 'USD',
    destinationCurrency: 'GBP',
    fxRate: 0.7891,
    convertedAmount: 39455,
    status: 'completed',
    createdAt: '2024-01-15T10:30:00Z',
    updatedAt: '2024-01-15T14:45:00Z',
  },
  {
    id: 'txn-002',
    reference: 'TXN-2024-001235',
    sender: mockAccounts[1],
    recipient: { name: 'Eurotrade GmbH', accountNumber: 'DE1234567890', country: 'Germany', currency: 'EUR' },
    amount: 25000,
    sourceCurrency: 'EUR',
    destinationCurrency: 'USD',
    fxRate: 1.0850,
    convertedAmount: 27125,
    status: 'processing',
    createdAt: '2024-01-15T11:00:00Z',
    updatedAt: '2024-01-15T11:15:00Z',
  },
  {
    id: 'txn-003',
    reference: 'TXN-2024-001236',
    sender: mockAccounts[0],
    recipient: { name: 'Nigeria Holdings Ltd', accountNumber: 'NGN1122334455', country: 'Nigeria', currency: 'NGN' },
    amount: 100000,
    sourceCurrency: 'USD',
    destinationCurrency: 'NGN',
    fxRate: 895.50,
    convertedAmount: 89550000,
    status: 'failed',
    createdAt: '2024-01-15T09:00:00Z',
    updatedAt: '2024-01-15T09:30:00Z',
  },
  {
    id: 'txn-004',
    reference: 'TXN-2024-001237',
    sender: mockAccounts[2],
    recipient: { name: 'Tokyo Imports K.K.', accountNumber: 'JP1234567890', country: 'Japan', currency: 'JPY' },
    amount: 75000,
    sourceCurrency: 'GBP',
    destinationCurrency: 'JPY',
    fxRate: 189.25,
    convertedAmount: 14193750,
    status: 'completed',
    createdAt: '2024-01-14T16:00:00Z',
    updatedAt: '2024-01-14T18:30:00Z',
  },
  {
    id: 'txn-005',
    reference: 'TXN-2024-001238',
    sender: mockAccounts[3],
    recipient: { name: 'Global Payments Inc', accountNumber: 'US9988776655', country: 'United States', currency: 'USD' },
    amount: 50000000,
    sourceCurrency: 'NGN',
    destinationCurrency: 'USD',
    fxRate: 0.00112,
    convertedAmount: 56000,
    status: 'processing',
    createdAt: '2024-01-15T12:00:00Z',
    updatedAt: '2024-01-15T12:10:00Z',
  },
  {
    id: 'txn-006',
    reference: 'TXN-2024-001239',
    sender: mockAccounts[0],
    recipient: { name: 'Sydney Trading Co', accountNumber: 'AU1234567890', country: 'Australia', currency: 'AUD' },
    amount: 35000,
    sourceCurrency: 'USD',
    destinationCurrency: 'AUD',
    fxRate: 1.5345,
    convertedAmount: 53707.50,
    status: 'completed',
    createdAt: '2024-01-13T08:00:00Z',
    updatedAt: '2024-01-13T10:15:00Z',
  },
  {
    id: 'txn-007',
    reference: 'TXN-2024-001240',
    sender: mockAccounts[1],
    recipient: { name: 'Toronto Services Ltd', accountNumber: 'CA1234567890', country: 'Canada', currency: 'CAD' },
    amount: 45000,
    sourceCurrency: 'EUR',
    destinationCurrency: 'CAD',
    fxRate: 1.4520,
    convertedAmount: 65340,
    status: 'failed',
    createdAt: '2024-01-12T14:30:00Z',
    updatedAt: '2024-01-12T15:00:00Z',
  },
  {
    id: 'txn-008',
    reference: 'TXN-2024-001241',
    sender: mockAccounts[2],
    recipient: { name: 'Paris Boutique SARL', accountNumber: 'FR1234567890', country: 'France', currency: 'EUR' },
    amount: 20000,
    sourceCurrency: 'GBP',
    destinationCurrency: 'EUR',
    fxRate: 1.1725,
    convertedAmount: 23450,
    status: 'completed',
    createdAt: '2024-01-11T09:45:00Z',
    updatedAt: '2024-01-11T11:20:00Z',
  },
];

// Mock transaction events (audit log)
export const mockTransactionEvents: Record<string, TransactionEvent[]> = {
  'txn-001': [
    { id: 'evt-001-1', transactionId: 'txn-001', status: 'processing', description: 'Transaction initiated', timestamp: '2024-01-15T10:30:00Z' },
    { id: 'evt-001-2', transactionId: 'txn-001', status: 'processing', description: 'FX conversion executed at rate 0.7891', timestamp: '2024-01-15T12:00:00Z' },
    { id: 'evt-001-3', transactionId: 'txn-001', status: 'completed', description: 'Payment sent to recipient bank', timestamp: '2024-01-15T14:00:00Z' },
    { id: 'evt-001-4', transactionId: 'txn-001', status: 'completed', description: 'Transaction completed - funds received by beneficiary', timestamp: '2024-01-15T14:45:00Z' },
  ],
  'txn-002': [
    { id: 'evt-002-1', transactionId: 'txn-002', status: 'processing', description: 'Transaction initiated', timestamp: '2024-01-15T11:00:00Z' },
    { id: 'evt-002-2', transactionId: 'txn-002', status: 'processing', description: 'FX conversion executed at rate 1.0850', timestamp: '2024-01-15T11:15:00Z' },
  ],
  'txn-003': [
    { id: 'evt-003-1', transactionId: 'txn-003', status: 'processing', description: 'Transaction initiated', timestamp: '2024-01-15T09:00:00Z' },
    { id: 'evt-003-2', transactionId: 'txn-003', status: 'failed', description: 'Payment rejected - invalid recipient account', timestamp: '2024-01-15T09:30:00Z' },
  ],
  'txn-005': [
    { id: 'evt-005-1', transactionId: 'txn-005', status: 'processing', description: 'Transaction initiated', timestamp: '2024-01-15T12:00:00Z' },
    { id: 'evt-005-2', transactionId: 'txn-005', status: 'processing', description: 'FX conversion executed at rate 0.00112', timestamp: '2024-01-15T12:10:00Z' },
  ],
};

// Mock ledger entries
export const mockLedgerEntries: Record<string, LedgerEntry[]> = {
  'txn-001': [
    { id: 'led-001-1', transactionId: 'txn-001', accountId: 'acc-001', accountName: 'Operating Account USD', type: 'debit', amount: 50000, currency: 'USD', balanceAfter: 1200000, timestamp: '2024-01-15T10:30:00Z' },
    { id: 'led-001-2', transactionId: 'txn-001', accountId: 'acc-001', accountName: 'Operating Account USD', type: 'credit', amount: 39455, currency: 'GBP', balanceAfter: 1200000, timestamp: '2024-01-15T14:45:00Z' },
  ],
  'txn-002': [
    { id: 'led-002-1', transactionId: 'txn-002', accountId: 'acc-002', accountName: 'EUR Treasury', type: 'debit', amount: 25000, currency: 'EUR', balanceAfter: 825000, timestamp: '2024-01-15T11:00:00Z' },
  ],
  'txn-003': [
    { id: 'led-003-1', transactionId: 'txn-003', accountId: 'acc-001', accountName: 'Operating Account USD', type: 'debit', amount: 100000, currency: 'USD', balanceAfter: 1150000, timestamp: '2024-01-15T09:00:00Z' },
    { id: 'led-003-2', transactionId: 'txn-003', accountId: 'acc-001', accountName: 'Operating Account USD', type: 'credit', amount: 100000, currency: 'USD', balanceAfter: 1250000, timestamp: '2024-01-15T09:30:00Z' },
  ],
};

// Simulated FX rate lookup
const fxRates: Record<string, number> = {
  'USD-GBP': 0.7891,
  'USD-EUR': 0.9215,
  'USD-NGN': 895.50,
  'USD-JPY': 148.25,
  'USD-CAD': 1.3520,
  'USD-AUD': 1.5345,
  'GBP-USD': 1.2675,
  'GBP-EUR': 1.1725,
  'GBP-NGN': 1135.00,
  'GBP-JPY': 189.25,
  'GBP-CAD': 1.7125,
  'GBP-AUD': 1.9435,
  'EUR-USD': 1.0850,
  'EUR-GBP': 0.8525,
  'EUR-NGN': 971.50,
  'EUR-JPY': 160.75,
  'EUR-CAD': 1.4520,
  'EUR-AUD': 1.6480,
  'NGN-USD': 0.00112,
  'NGN-GBP': 0.00088,
  'NGN-EUR': 0.00103,
  'NGN-JPY': 0.1655,
  'NGN-CAD': 0.00151,
  'NGN-AUD': 0.00171,
};

export function getFXRate(from: Currency, to: Currency): FXQuote {
  const key = `${from}-${to}`;
  const rate = fxRates[key] || 1;
  
  return {
    fromCurrency: from,
    toCurrency: to,
    rate,
    validUntil: new Date(Date.now() + 5 * 60 * 1000).toISOString(), // 5 minutes validity
  };
}

export function getAccounts(): Account[] {
  return mockAccounts;
}

export function getTransaction(id: string): Transaction | undefined {
  return mockTransactions.find(t => t.id === id);
}

export function getTransactions(
  statusFilter?: TransactionStatus | 'all',
  page: number = 1,
  limit: number = 10
): { transactions: Transaction[]; total: number; totalPages: number } {
  let filtered = [...mockTransactions];
  
  if (statusFilter && statusFilter !== 'all') {
    filtered = filtered.filter(t => t.status === statusFilter);
  }
  
  const total = filtered.length;
  const totalPages = Math.ceil(total / limit);
  const start = (page - 1) * limit;
  const transactions = filtered.slice(start, start + limit);
  
  return { transactions, total, totalPages };
}

export function getTransactionEvents(transactionId: string): TransactionEvent[] {
  return mockTransactionEvents[transactionId] || [];
}

export function getLedgerEntries(transactionId: string): LedgerEntry[] {
  return mockLedgerEntries[transactionId] || [];
}

export async function createTransaction(data: {
  sourceAccountId: string;
  recipientName: string;
  recipientAccountNumber: string;
  recipientCountry: string;
  destinationCurrency: Currency;
  amount: number;
}): Promise<{ success: boolean; transaction?: Transaction; error?: string }> {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1500));
  
  const sourceAccount = mockAccounts.find(a => a.id === data.sourceAccountId);
  if (!sourceAccount) {
    return { success: false, error: 'Invalid source account' };
  }
  
  if (data.amount > sourceAccount.balance) {
    return { success: false, error: 'Insufficient funds in source account' };
  }
  
  const fxQuote = getFXRate(sourceAccount.currency, data.destinationCurrency);
  const convertedAmount = data.amount * fxQuote.rate;
  
  const newTransaction: Transaction = {
    id: `txn-${Date.now()}`,
    reference: `TXN-${new Date().getFullYear()}-${String(mockTransactions.length + 1).padStart(6, '0')}`,
    sender: sourceAccount,
    recipient: {
      name: data.recipientName,
      accountNumber: data.recipientAccountNumber,
      country: data.recipientCountry,
      currency: data.destinationCurrency,
    },
    amount: data.amount,
    sourceCurrency: sourceAccount.currency,
    destinationCurrency: data.destinationCurrency,
    fxRate: fxQuote.rate,
    convertedAmount,
    status: 'processing',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };
  
  mockTransactions.unshift(newTransaction);
  mockTransactionEvents[newTransaction.id] = [
    {
      id: `evt-${Date.now()}`,
      transactionId: newTransaction.id,
      status: 'processing',
      description: 'Transaction initiated',
      timestamp: new Date().toISOString(),
    },
  ];
  
  return { success: true, transaction: newTransaction };
}
