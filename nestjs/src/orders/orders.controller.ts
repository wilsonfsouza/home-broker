import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { InputExecuteTransactionDTO, InputTransactionDTO } from './order.dto';
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
  executeTransaction(@Body() body: InputExecuteTransactionDTO) {
    this.ordersService.executeTransaction(body);
  }
}
