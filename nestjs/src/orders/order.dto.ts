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
