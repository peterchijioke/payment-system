import { NextRequest, NextResponse } from 'next/server';
import { paymentApi } from '@/lib/api/server-api';

/**
 * Server-side API route that proxies requests to the backend
 * The backend URL is kept secure on the server and never exposed to clients
 * 
 * This route handles:
 * - GET /api/payments - List all payments
 */

export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const limit = searchParams.get('limit');
    const offset = searchParams.get('offset');
    const status = searchParams.get('status');
    const start_date = searchParams.get('start_date');
    const end_date = searchParams.get('end_date');

    const params: {
      limit?: number;
      offset?: number;
      status?: string;
      start_date?: string;
      end_date?: string;
    } = {};
    
    if (limit) params.limit = parseInt(limit, 10);
    if (offset) params.offset = parseInt(offset, 10);
    if (status) params.status = status;
    if (start_date) params.start_date = start_date;
    if (end_date) params.end_date = end_date;

    const response = await paymentApi.getPayments(params);
    
    return NextResponse.json(response);
  } catch (error: any) {
    console.error('Error fetching payments:', error);
    
    return NextResponse.json(
      { success: false, error: 'Failed to fetch payments' },
      { status: error.response?.status || 500 }
    );
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    // Extract idempotency key from request headers or generate one
    const idempotencyKey = request.headers.get('Idempotency-Key') || `payment-${Date.now()}`;
    
    // Remove idempotency_key from body if present (we pass it via header)
    const { idempotency_key, ...paymentData } = body;
    
    const response = await paymentApi.createPayment({
      ...paymentData,
      idempotency_key: idempotency_key || idempotencyKey,
    });
    
    // Create response with cookie containing idempotency key
    const nextResponse = NextResponse.json(response);
    
    // Set cookie with idempotency key for future requests
    nextResponse.cookies.set('last_idempotency_key', idempotencyKey, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 60 * 60 * 24, // 24 hours
      path: '/',
    });
    
    return nextResponse;
  } catch (error: any) {
    console.error('Error creating payment:', error);
    
    return NextResponse.json(
      { success: false, error: 'Failed to create payment' },
      { status: error.response?.status || 500 }
    );
  }
}
