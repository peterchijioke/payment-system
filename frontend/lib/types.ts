// Transaction types and interfaces

export type TransactionStatus = 'processing' | 'completed' | 'failed';

export type Currency = 'USD' | 'EUR' | 'GBP' | 'NGN' | 'JPY' | 'CAD' | 'AUD';

export interface Account {
  id: string;
  name: string;
  accountNumber: string;
  currency: Currency;
  balance: number;
}

export interface Transaction {
  id: string;
  reference: string;
  sender: Account;
  recipient: {
    name: string;
    accountNumber: string;
    country: string;
    currency: Currency;
  };
  amount: number;
  sourceCurrency: Currency;
  destinationCurrency: Currency;
  fxRate: number;
  convertedAmount: number;
  status: TransactionStatus;
  createdAt: string;
  updatedAt: string;
}

export interface TransactionEvent {
  id: string;
  transactionId: string;
  status: TransactionStatus;
  description: string;
  timestamp: string;
}

export interface LedgerEntry {
  id: string;
  transactionId: string;
  accountId: string;
  accountName: string;
  type: 'debit' | 'credit';
  amount: number;
  currency: Currency;
  balanceAfter: number;
  timestamp: string;
}

export interface FXQuote {
  fromCurrency: Currency;
  toCurrency: Currency;
  rate: number;
  validUntil: string;
}

export interface PaymentFormData {
  sourceAccountId: string;
  recipientName: string;
  recipientAccountNumber: string;
  recipientCountry: string;
  destinationCurrency: Currency;
  amount: number;
}
