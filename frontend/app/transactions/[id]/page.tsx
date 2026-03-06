'use client';

import { use } from 'react';
import { TransactionDetail } from '@/components/transaction-detail';

export default function TransactionDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  
  return (
    <div className="min-h-screen bg-zinc-50 p-6">
      <div className="max-w-7xl mx-auto">
        <TransactionDetail id={id} />
      </div>
    </div>
  );
}
