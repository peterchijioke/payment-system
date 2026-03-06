'use client';

import { useState } from 'react';
import Link from 'next/link';
import { format } from 'date-fns';
import { useQuery } from '@tanstack/react-query';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { 
  Pagination, 
  PaginationContent, 
  PaginationItem, 
  PaginationLink, 
  PaginationNext, 
  PaginationPrevious 
} from '@/components/ui/pagination';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface TransactionListProps {
  onTransactionSelect?: (transaction: any) => void;
}

interface PaymentTransaction {
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
}

interface Payment {
  transaction: PaymentTransaction;
  ledger_entries: any[];
  timeline: { status: string; timestamp: string }[];
}

interface PaymentsResponse {
  data: {
    limit: number;
    offset: number;
    payments: Payment[];
    total: number;
  };
  success: boolean;
}

type TransactionStatus = 'initiated' | 'processing' | 'settled' | 'completed' | 'failed' | 'reversed' | 'pending_review';

const statusConfig: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
  initiated: { label: 'Initiated', variant: 'outline' },
  processing: { label: 'Processing', variant: 'secondary' },
  settled: { label: 'Settled', variant: 'secondary' },
  completed: { label: 'Completed', variant: 'default' },
  failed: { label: 'Failed', variant: 'destructive' },
  reversed: { label: 'Reversed', variant: 'outline' },
  pending_review: { label: 'Pending Review', variant: 'outline' },
};

function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

async function fetchPayments(params: {
  status?: string;
  limit: number;
  offset: number;
}): Promise<PaymentsResponse> {
  const searchParams = new URLSearchParams();
  if (params.status && params.status !== 'all') {
    searchParams.set('status', params.status);
  }
  searchParams.set('limit', params.limit.toString());
  searchParams.set('offset', params.offset.toString());

  const response = await fetch(`/api/payments?${searchParams.toString()}`);
  if (!response.ok) {
    throw new Error('Failed to fetch payments');
  }
  return response.json();
}

export function TransactionList({ onTransactionSelect }: TransactionListProps) {
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error } = useQuery<PaymentsResponse>({
    queryKey: ['payments', statusFilter, currentPage],
    queryFn: () => fetchPayments({
      status: statusFilter,
      limit: pageSize,
      offset: (currentPage - 1) * pageSize,
    }),
  });

  const transactions = data?.data?.payments || [];
  const total = data?.data?.total || 0;
  const totalPages = Math.ceil(total / pageSize);

  const handleStatusFilterChange = (value: string) => {
    setStatusFilter(value);
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle>Transactions</CardTitle>
        <div className="flex items-center gap-4">
          <Select value={statusFilter} onValueChange={handleStatusFilterChange}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Statuses</SelectItem>
              <SelectItem value="initiated">Initiated</SelectItem>
              <SelectItem value="processing">Processing</SelectItem>
              <SelectItem value="settled">Settled</SelectItem>
              <SelectItem value="completed">Completed</SelectItem>
              <SelectItem value="failed">Failed</SelectItem>
              <SelectItem value="reversed">Reversed</SelectItem>
              <SelectItem value="pending_review">Pending Review</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading transactions...</div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center h-64">
            <div className="text-red-500">Error loading transactions</div>
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Reference</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>FX Rate</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {transactions.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      No transactions found.
                    </TableCell>
                  </TableRow>
                ) : (
                  transactions.map((payment) => (
                    <TableRow 
                      key={payment.transaction.id}
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => onTransactionSelect?.(payment)}
                    >
                      <TableCell className="font-medium">
                        <Link 
                          href={`/transactions/${payment.transaction.id}`}
                          className="hover:underline text-primary"
                        >
                          {payment.transaction.transaction_reference}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <div>
                          <div className="font-medium">{payment.transaction.description || '-'}</div>
                          <div className="text-sm text-muted-foreground">
                            {payment.transaction.reference || 'No reference'}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="space-y-1">
                          <div>{formatCurrency(payment.transaction.amount, payment.transaction.currency)}</div>
                          <div className="text-sm text-muted-foreground">
                            → {formatCurrency(payment.transaction.fx_amount, payment.transaction.fx_currency)}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="font-mono text-sm">
                        {payment.transaction.fx_rate.toFixed(6)}
                      </TableCell>
                      <TableCell>
                        <Badge variant={statusConfig[payment.transaction.status]?.variant || 'outline'}>
                          {statusConfig[payment.transaction.status]?.label || payment.transaction.status}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="text-sm">
                          {format(new Date(payment.transaction.initiated_at), 'MMM d, yyyy')}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {format(new Date(payment.transaction.initiated_at), 'HH:mm')}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
            
            {totalPages > 1 && (
              <div className="mt-4 flex items-center justify-between">
                <div className="text-sm text-muted-foreground">
                  Showing {((currentPage - 1) * pageSize) + 1} to {Math.min(currentPage * pageSize, total)} of {total} transactions
                </div>
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious 
                        onClick={() => handlePageChange(currentPage - 1)}
                        className={currentPage === 1 ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                      />
                    </PaginationItem>
                    {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                      <PaginationItem key={page}>
                        <PaginationLink 
                          onClick={() => handlePageChange(page)}
                          isActive={currentPage === page}
                          className="cursor-pointer"
                        >
                          {page}
                        </PaginationLink>
                      </PaginationItem>
                    ))}
                    <PaginationItem>
                      <PaginationNext 
                        onClick={() => handlePageChange(currentPage + 1)}
                        className={currentPage === totalPages ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                      />
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
