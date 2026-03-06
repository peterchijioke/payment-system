import { NextRequest, NextResponse } from 'next/server';
import { paymentApi } from '@/lib/api/server-api';

/**
 * Server-side API route for getting a single payment by ID
 * The backend URL is kept secure on the server
 */

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const response = await paymentApi.getPayment(id);
    
    // Handle payment not found
    if (!response.success && response.error) {
      return NextResponse.json(
        { success: false, error: response.error },
        { status: 404 }
      );
    }
    
    return NextResponse.json(response);
  } catch (error: any) {
    console.error('Error fetching payment:', error);
    
    // Handle not found error from backend
    if (error.response?.status === 404 || error.response?.data?.error?.includes('not found')) {
      return NextResponse.json(
        { success: false, error: 'Payment not found' },
        { status: 404 }
      );
    }
    
    return NextResponse.json(
      { success: false, error: 'Failed to fetch payment' },
      { status: error.response?.status || 500 }
    );
  }
}
