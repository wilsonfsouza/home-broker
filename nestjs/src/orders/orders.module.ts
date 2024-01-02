import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { OrdersController } from './orders.controller';
import { OrdersService } from './orders.service';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'ORDERS_PUBLISHER',
        transport: Transport.KAFKA,
        options: {
          client: {
            clientId: 'orders',
            // brokers: ['host.docker.internal:9094'], // docker default host
            brokers: ['kafka:9092'], // docker internal network
          },
        },
      },
    ]),
  ],
  controllers: [OrdersController],
  providers: [OrdersService],
})
export class OrdersModule {}
