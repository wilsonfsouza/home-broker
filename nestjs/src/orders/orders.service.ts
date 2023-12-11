import { Injectable } from '@nestjs/common';
import { OrderStatus } from '@prisma/client';
import { PrismaService } from 'src/prisma/prisma/prisma.service';
import { InputExecuteTransactionDTO, InputTransactionDTO } from './order.dto';

@Injectable()
export class OrdersService {
  constructor(private prismaService: PrismaService) {}

  all(filter: { wallet_id: string }) {
    return this.prismaService.order.findMany({
      where: {
        wallet_id: filter.wallet_id,
      },
      include: {
        Transactions: true,
        Asset: {
          select: {
            id: true,
            symbol: true,
          },
        },
      },
      orderBy: {
        updated_at: 'desc',
      },
    });
  }

  startTransaction(data: InputTransactionDTO) {
    return this.prismaService.order.create({
      data: {
        asset_id: data.asset_id,
        wallet_id: data.wallet_id,
        shares: data.shares,
        price: data.price,
        type: data.type,
        status: OrderStatus.PENDING,
        partial: data.shares,
      },
    });
  }

  async executeTransaction(data: InputExecuteTransactionDTO) {
    return this.prismaService.$transaction(async (prisma) => {
      const order = await prisma.order.findUniqueOrThrow({
        where: { id: data.order_id },
      });

      await prisma.order.update({
        where: {
          id: data.order_id,
        },
        data: {
          status: data.status,
          partial: order.partial - data.negotiated_shares,
          Transactions: {
            create: {
              broker_transaction_id: data.broker_transaction_id,
              related_investor_id: data.related_investor_id,
              shares: data.negotiated_shares,
              price: data.price,
            },
          },
        },
      });

      if (data.status === OrderStatus.CLOSED) {
        await prisma.asset.update({
          where: { id: order.asset_id },
          data: {
            price: data.price,
          },
        });

        const walletAsset = await prisma.walletAsset.findUnique({
          where: {
            wallet_id_asset_id: {
              asset_id: order.asset_id,
              wallet_id: order.wallet_id,
            },
          },
        });

        if (!walletAsset && order.type === 'BUY') {
          await prisma.walletAsset.create({
            data: {
              asset_id: order.asset_id,
              wallet_id: order.wallet_id,
              shares: data.negotiated_shares,
            },
          });
        }

        if (walletAsset) {
          await prisma.walletAsset.update({
            where: {
              wallet_id_asset_id: {
                asset_id: order.asset_id,
                wallet_id: order.wallet_id,
              },
            },
            data: {
              shares:
                order.type === 'BUY'
                  ? walletAsset.shares + order.shares
                  : walletAsset.shares - order.shares,
            },
          });
        }
      }
    });
  }
}
