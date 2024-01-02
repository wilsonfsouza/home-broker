import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import {
  ExecuteTransactionMessageDTO,
  InputExecuteTransactionDTO,
  InputTransactionDTO,
} from './order.dto';
import { OrdersService } from './orders.service';

@Controller('wallets/:wallet_id/orders')
export class OrdersController {
  constructor(private readonly ordersService: OrdersService) {}

  @Get()
  all(@Param('wallet_id') wallet_id: string) {
    return this.ordersService.all({ wallet_id });
  }

  @Post()
  startTransaction(
    @Param('wallet_id') wallet_id: string,
    @Body() body: Omit<InputTransactionDTO, 'wallet_id'>,
  ) {
    return this.ordersService.startTransaction({
      ...body,
      wallet_id,
    });
  }

  @Post('execute')
  executeTransactionRest(@Body() body: InputExecuteTransactionDTO) {
    this.ordersService.executeTransaction(body);
  }

  @MessagePattern('output') // output is the topic of the kafka message
  async executeTransactionConsumer(
    @Payload() message: ExecuteTransactionMessageDTO,
  ) {
    const transaction = message.transactions[message.transactions.length - 1];

    await this.ordersService.executeTransaction({
      order_id: message.order_id,
      status: message.status,
      related_investor_id:
        message.order_type === 'BUY'
          ? transaction.seller_id
          : transaction.buyer_id,
      broker_transaction_id: transaction.transaction_id,
      negotiated_shares: transaction.shares,
      price: transaction.price,
    });
  }
}
