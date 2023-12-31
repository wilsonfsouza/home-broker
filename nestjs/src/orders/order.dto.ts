import { OrderStatus, OrderType } from '@prisma/client';

export class InputTransactionDTO {
  asset_id: string;
  wallet_id: string;
  shares: number;
  price: number;
  type: OrderType;
}

export class InputExecuteTransactionDTO {
  order_id: string;
  status: OrderStatus;
  price: number;
  related_investor_id: string;
  broker_transaction_id: string;
  negotiated_shares: number;
}

type TransctionMessageOutput = {
  transaction_id: string;
  buyer_id: string;
  seller_id: string;
  asset_id: string;
  shares: number;
  price: number;
};

export class ExecuteTransactionMessageDTO {
  order_id: string;
  investor_id: string;
  asset_id: string;
  order_type: string;
  status: 'OPEN' | 'CLOSED';
  partial: number;
  shares: number;
  transactions: TransctionMessageOutput[];
}
