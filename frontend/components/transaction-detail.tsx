'use client';

import { use } from 'react';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { format } from 'date-fns';
import { ArrowLeft, Clock, CheckCircle2, XCircle, Loader2, Building2, User, ArrowRightLeft } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';

interface TransactionDetailProps {
  id: string;
}

interface TransactionData {
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
  ledger_entries: {
    id: string;
    entry_reference: string;
    transaction_id: string;
    account_id: string;
    entry_type: string;
    amount: number;
    currency: string;
    counterpart_entry_id: string | null;
    original_entry_id: string | null;
    status: string;
    reversal_reason: string;
    description: string;
    effective_date: string;
    posted_at: string | null;
    reversed_by_id: string | null;
    created_at: string;
    created_by: string | null;
  }[];
  timeline: {
    status: string;
    timestamp: string;
  }[];
}

const statusConfig: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; icon: typeof CheckCircle2 }> = {
  initiated: { label: 'Initiated', variant: 'outline', icon: Clock },
  processing: { label: 'Processing', variant: 'secondary', icon: Loader2 },
  settled: { label: 'Settled', variant: 'secondary', icon: Loader2 },
  completed: { label: 'Completed', variant: 'default', icon: CheckCircle2 },
  failed: { label: 'Failed', variant: 'destructive', icon: XCircle },
  reversed: { label: 'Reversed', variant: 'outline', icon: XCircle },
  pending_review: { label: 'Pending Review', variant: 'outline', icon: Clock },
};

function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

async function fetchPayment(id: string): Promise<{ data: TransactionData; success: boolean }> {
  const response = await fetch(`/api/payments/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch payment');
  }
  return response.json();
}

export function TransactionDetail({ id }: TransactionDetailProps) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['payment', id],
    queryFn: () => fetchPayment(id),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !data?.success) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <p className="text-muted-foreground mb-4">Transaction not found</p>
          <Button asChild>
            <Link href="/">Back to Dashboard</Link>
          </Button>
        </CardContent>
      </Card>
    );
  }

  const transaction = data.data.transaction;
  const ledgerEntries = data.data.ledger_entries;
  const timeline = data.data.timeline;

  const statusInfo = statusConfig[transaction.status] || { label: transaction.status, variant: 'outline' as const, icon: Clock };
  const StatusIcon = statusInfo.icon;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Transaction Details</h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Transaction Overview */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-3">
              {transaction.transaction_reference}
              <Badge variant={statusInfo.variant} className="flex items-center gap-1">
                <StatusIcon className={`h-3 w-3 ${transaction.status === 'processing' || transaction.status === 'settled' ? 'animate-spin' : ''}`} />
                {statusInfo.label}
              </Badge>
            </CardTitle>
            <CardDescription>
              Created {format(new Date(transaction.initiated_at), 'MMMM d, yyyy \'at\' HH:mm')}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Amount Details */}
            <div className="border rounded-lg p-4 bg-muted/50">
              <div className="flex items-center gap-2 mb-4">
                <ArrowRightLeft className="h-4 w-4" />
                <span className="font-semibold">Payment Details</span>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-muted-foreground">Source Amount</p>
                  <p className="text-lg font-medium">
                    {formatCurrency(transaction.amount, transaction.currency)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Destination Amount</p>
                  <p className="text-lg font-medium">
                    {formatCurrency(transaction.fx_amount, transaction.fx_currency)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Exchange Rate</p>
                  <p className="font-mono">
                    1 {transaction.currency} = {transaction.fx_rate.toFixed(6)} {transaction.fx_currency}
                  </p>
                </div>
              </div>
            </div>

            {/* Description */}
            {transaction.description && (
              <div className="space-y-2">
                <h3 className="font-semibold flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Description
                </h3>
                <div className="bg-muted/30 rounded-lg p-4">
                  <p className="font-medium">{transaction.description}</p>
                </div>
              </div>
            )}

            {/* Reference */}
            {transaction.reference && (
              <div className="grid grid-cols-2 gap-4 pt-4 border-t">
                <div>
                  <p className="text-sm text-muted-foreground">Reference</p>
                  <p className="font-medium">{transaction.reference}</p>
                </div>
              </div>
            )}

            {/* Timestamps */}
            <div className="grid grid-cols-2 gap-4 pt-4 border-t">
              <div>
                <p className="text-sm text-muted-foreground">Initiated</p>
                <p>{format(new Date(transaction.initiated_at), 'MMM d, yyyy HH:mm:ss')}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Last Updated</p>
                <p>{format(new Date(transaction.updated_at), 'MMM d, yyyy HH:mm:ss')}</p>
              </div>
              {transaction.settled_at && (
                <div>
                  <p className="text-sm text-muted-foreground">Settled</p>
                  <p>{format(new Date(transaction.settled_at), 'MMM d, yyyy HH:mm:ss')}</p>
                </div>
              )}
              {transaction.completed_at && (
                <div>
                  <p className="text-sm text-muted-foreground">Completed</p>
                  <p>{format(new Date(transaction.completed_at), 'MMM d, yyyy HH:mm:ss')}</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Timeline / Audit Log */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-4 w-4" />
              Timeline
            </CardTitle>
          </CardHeader>
          <CardContent>
            {timeline.length === 0 ? (
              <p className="text-muted-foreground text-sm">No events recorded</p>
            ) : (
              <div className="relative">
                {/* Vertical line */}
                <div className="absolute left-3 top-0 bottom-0 w-0.5 bg-border" />
                
                <div className="space-y-6">
                  {timeline.map((event, index) => {
                    const eventStatusConfig = statusConfig[event.status] || { label: event.status, variant: 'outline' as const, icon: Clock };
                    const EventIcon = eventStatusConfig.icon;
                    
                    return (
                      <div key={index} className="relative pl-8">
                        {/* Dot */}
                        <div className={`absolute left-1.5 top-1 w-3 h-3 rounded-full border-2 border-background ${
                          event.status === 'completed' ? 'bg-green-500' :
                          event.status === 'failed' ? 'bg-destructive' :
                          event.status === 'settled' ? 'bg-blue-500' :
                          'bg-yellow-500'
                        }`} />
                        
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <Badge variant={eventStatusConfig.variant} className="text-xs">
                              <EventIcon className={`h-3 w-3 mr-1 ${event.status === 'processing' || event.status === 'settled' ? 'animate-spin' : ''}`} />
                              {eventStatusConfig.label}
                            </Badge>
                          </div>
                          <p className="text-xs text-muted-foreground">
                            {format(new Date(event.timestamp), 'MMM d, HH:mm:ss')}
                          </p>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Ledger Entries */}
      <Card>
        <CardHeader>
          <CardTitle>Ledger Entries</CardTitle>
          <CardDescription>
            Associated ledger entries for this transaction
          </CardDescription>
        </CardHeader>
        <CardContent>
          {ledgerEntries.length === 0 ? (
            <p className="text-muted-foreground">No ledger entries found</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Entry ID</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ledgerEntries.map((entry) => (
                  <TableRow key={entry.id}>
                    <TableCell className="font-mono text-sm">{entry.entry_reference}</TableCell>
                    <TableCell>
                      <Badge variant={entry.entry_type === 'debit' ? 'destructive' : 'default'}>
                        {entry.entry_type.toUpperCase()}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono">
                      {formatCurrency(entry.amount, entry.currency)}
                    </TableCell>
                    <TableCell>
                      <Badge variant={entry.status === 'posted' ? 'default' : 'secondary'}>
                        {entry.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm">
                      {entry.description}
                    </TableCell>
                    <TableCell>
                      {format(new Date(entry.created_at), 'MMM d, yyyy HH:mm:ss')}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
