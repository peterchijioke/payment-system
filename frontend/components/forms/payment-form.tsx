'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { format } from 'date-fns';
import { z } from 'zod';
import { Loader2, ArrowRightLeft, AlertCircle, CheckCircle2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { countries } from '@/lib/countries';

interface PaymentFormProps {
  onSuccess?: () => void;
}

const currencies = ['USD', 'EUR', 'GBP', 'NGN', 'JPY', 'CAD', 'AUD'];

const paymentSchema = z.object({
  sourceAccountId: z.string().min(1, 'Please select a source account'),
  recipientName: z.string().min(1, 'Recipient name is required').min(2, 'Name must be at least 2 characters'),
  recipientAccountNumber: z.string().min(1, 'Account number is required').min(4, 'Account number must be at least 4 digits'),
  recipientCountry: z.string().min(1, 'Country is required').min(2, 'Country must be at least 2 characters'),
  recipientBank: z.string().optional(),
  destinationCurrency: z.string().min(1, 'Please select a destination currency'),
  amount: z.coerce.number().min(1, 'Amount must be greater than 0'),
});

type PaymentFormValues = z.infer<typeof paymentSchema>;


const mockAccounts: Record<string, { name: string; accountNumber: string; currency: string; balance: number }> = {
  'd5c54768-3c8d-472d-8eb8-0a6871ff6147': {
    name: 'Operating Account',
    accountNumber: 'US1234567890',
    currency: 'USD',
    balance: 1000000
  }
};

function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

async function createPayment(data: {
  account_id: string;
  amount: number;
  currency: string;
  destination_currency: string;
  recipient_name: string;
  recipient_account: string;
  recipient_bank: string;
  recipient_country: string;
  reference?: string;
}) {
  const cookies = document.cookie.split('; ').find(c => c.startsWith('last_idempotency_key='));
  const existingKey = cookies?.split('=')[1];
  
  const idempotencyKey = existingKey || `payment-${Date.now()}-${Math.random().toString(36).substring(7)}`;
  
  const response = await fetch('/api/payments', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Idempotency-Key': idempotencyKey,
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create payment');
  }
  
  return response.json();
}

export function PaymentForm({ onSuccess }: PaymentFormProps) {
  const queryClient = useQueryClient();
  const [fxRate, setFxRate] = useState<number | null>(null);
  
  // Get accounts as array
  const accounts = Object.entries(mockAccounts).map(([id, details]) => ({
    id,
    ...details,
  }));

  // Initialize form with Zod resolver
  const form = useForm<PaymentFormValues>({
    resolver: zodResolver(paymentSchema),
    defaultValues: {
      sourceAccountId: '',
      recipientName: '',
      recipientAccountNumber: '',
      recipientCountry: '',
      recipientBank: '',
      destinationCurrency: '',
      amount: 0,
    },
    mode: 'onChange',
  });

  // Watch form values for FX rate calculation
  const watchedValues = form.watch();

  // React Query mutation for payment creation
  const mutation = useMutation({
    mutationFn: createPayment,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      onSuccess?.();
      // Reset form fields after successful payment
      form.reset();
      setFxRate(null);
    },
  });

  useEffect(() => {
    const sourceCurrency = mockAccounts[watchedValues.sourceAccountId]?.currency || 'USD';
    const destCurrency = watchedValues.destinationCurrency;
    
    if (sourceCurrency && destCurrency && sourceCurrency !== destCurrency && watchedValues.amount > 0) {
      const rates: Record<string, number> = {
        'USD-NGN': 895.50,
        'USD-EUR': 0.92,
        'USD-GBP': 0.79,
        'EUR-USD': 1.09,
        'EUR-NGN': 975.00,
        'GBP-USD': 1.27,
        'GBP-NGN': 1135.00,
        'NGN-USD': 0.00112,
        'NGN-EUR': 0.00103,
        'NGN-GBP': 0.00088,
      };
      setFxRate(rates[`${sourceCurrency}-${destCurrency}`] || null);
    } else {
      setFxRate(null);
    }
  }, [watchedValues.sourceAccountId, watchedValues.destinationCurrency, watchedValues.amount]);

  const sourceCurrency = mockAccounts[watchedValues.sourceAccountId]?.currency || 'USD';
  const destCurrency = watchedValues.destinationCurrency;

  const onSubmit = (data: PaymentFormValues) => {
    const account = mockAccounts[data.sourceAccountId];
    
    mutation.mutate({
      account_id: data.sourceAccountId,
      amount: data.amount,
      currency: account?.currency || 'USD',
      destination_currency: data.destinationCurrency,
      recipient_name: data.recipientName,
      recipient_account: data.recipientAccountNumber,
      recipient_bank: data.recipientBank || 'Bank',
      recipient_country: data.recipientCountry,
      reference: `PAY-${Date.now()}`,
    });
  };

  const sourceAccount = accounts.find(a => a.id === watchedValues.sourceAccountId);
  const convertedAmount = fxRate && watchedValues.amount ? watchedValues.amount * fxRate : 0;
  const amountExceedsBalance = sourceAccount && watchedValues.amount > 0 && watchedValues.amount > sourceAccount.balance;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Initiate Payment</CardTitle>
        <CardDescription>
          Create a new cross-border payment transaction
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {mutation.isSuccess && (
              <Alert className="mb-6 border-green-500 bg-green-50">
                <CheckCircle2 className="h-4 w-4 text-green-600" />
                <AlertTitle>Payment Initiated</AlertTitle>
                <AlertDescription>
                  Your payment has been successfully submitted. Transaction: {mutation.data?.data?.transaction_reference}
                </AlertDescription>
              </Alert>
            )}

            {mutation.isError && (
              <Alert variant="destructive" className="mb-6">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                  {(mutation.error as Error)?.message || 'Failed to create payment. Please try again.'}
                </AlertDescription>
              </Alert>
            )}

            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Source Account</h3>
              
              <FormField
                control={form.control}
                name="sourceAccountId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Account</FormLabel>
                    <Select 
                      onValueChange={field.onChange} 
                      defaultValue={field.value}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select source account" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {accounts.map((account) => (
                          <SelectItem key={account.id} value={account.id}>
                            <div className="flex items-center justify-between w-full">
                              <span>{account.name}</span>
                              <span className="text-muted-foreground ml-2">
                                ({account.currency}) - {formatCurrency(account.balance, account.currency)}
                              </span>
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormDescription>
                      Available balance: {sourceAccount 
                        ? formatCurrency(sourceAccount.balance, sourceAccount.currency)
                        : '--'}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Destination Details</h3>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="recipientName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Recipient Name</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter recipient name" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="recipientAccountNumber"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Recipient Account Number</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter account number" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="recipientCountry"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Country</FormLabel>
                      <Select 
                        onValueChange={field.onChange} 
                        defaultValue={field.value}
                        value={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select country" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {countries.map((country) => (
                            <SelectItem key={country.country} value={country.country}>
                              {country.country}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="recipientBank"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Bank Name (Optional)</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter bank name" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="destinationCurrency"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Destination Currency</FormLabel>
                    <Select 
                      onValueChange={field.onChange} 
                      defaultValue={field.value}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select currency" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {currencies.map((currency) => (
                          <SelectItem key={currency} value={currency}>
                            {currency}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Amount</h3>
              
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Amount ({sourceAccount?.currency || 'USD'})</FormLabel>
                    <FormControl>
                      <Input 
                        type="number" 
                        placeholder="Enter amount" 
                        {...field}
                      />
                    </FormControl>
                    {amountExceedsBalance && (
                      <FormDescription className="text-destructive">
                        Amount exceeds available balance
                      </FormDescription>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {fxRate && watchedValues.amount > 0 && sourceAccount && (
              <div className="border rounded-lg p-4 bg-muted/50">
                <div className="flex items-center gap-2 mb-4">
                  <ArrowRightLeft className="h-4 w-4" />
                  <span className="font-semibold">FX Quote</span>
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Exchange Rate</span>
                    <span className="font-mono">
                      1 {sourceCurrency} = {fxRate.toFixed(4)} {destCurrency}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">You Send</span>
                    <span className="font-medium">
                      {formatCurrency(watchedValues.amount, sourceCurrency)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Recipient Receives</span>
                    <span className="font-medium text-lg">
                      {formatCurrency(convertedAmount, destCurrency)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center pt-2 border-t">
                    <span className="text-sm text-muted-foreground">Quote Valid Until</span>
                    <span className="text-sm">
                      {format(new Date(Date.now() + 5 * 60 * 1000), 'HH:mm:ss')}
                    </span>
                  </div>
                </div>
              </div>
            )}

            <Button 
              type="submit" 
              className="w-full" 
              disabled={mutation.isPending || amountExceedsBalance}
            >
              {mutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Processing...
                </>
              ) : (
                'Submit Payment'
              )}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
