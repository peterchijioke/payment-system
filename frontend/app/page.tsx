'use client';

import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { TransactionList } from '@/components/transaction-list';
import { PaymentForm } from '@/components/forms/payment-form';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { List, PlusCircle, DollarSign } from 'lucide-react';

export default function Dashboard() {
  const [refreshKey, setRefreshKey] = useState(0);

  const handlePaymentSuccess = () => {
    // Refresh the transaction list
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-zinc-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-zinc-900">Operations Dashboard</h1>
            <p className="text-zinc-600 mt-1">Manage cross-border payments and monitor transactions</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Transactions</CardTitle>
              <List className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">8</div>
              <p className="text-xs text-muted-foreground">Across all statuses</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Processing</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">2</div>
              <p className="text-xs text-muted-foreground">Awaiting completion</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Failed</CardTitle>
              <PlusCircle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">2</div>
              <p className="text-xs text-muted-foreground">Requires attention</p>
            </CardContent>
          </Card>
        </div>

        <Tabs defaultValue="transactions" className="space-y-4">
          <TabsList>
            <TabsTrigger value="transactions" className="gap-2">
              <List className="h-4 w-4" />
              Transaction List
            </TabsTrigger>
            <TabsTrigger value="initiate" className="gap-2">
              <PlusCircle className="h-4 w-4" />
              Initiate Payment
            </TabsTrigger>
          </TabsList>
          
          <TabsContent value="transactions" className="space-y-4">
            <TransactionList key={refreshKey} />
          </TabsContent>
          
          <TabsContent value="initiate" className="space-y-4">
            <div className="max-w-2xl">
              <PaymentForm onSuccess={handlePaymentSuccess} />
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
