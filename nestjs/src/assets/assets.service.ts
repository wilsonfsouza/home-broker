import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma/prisma.service';

interface Asset {
  id: string;
  symbol: string;
  price: number;
}

@Injectable()
export class AssetsService {
  constructor(private prismaService: PrismaService) {}

  all() {
    return this.prismaService.asset.findMany();
  }

  create({ id, price, symbol }: Asset) {
    return this.prismaService.asset.create({
      data: {
        id,
        symbol,
        price,
      },
    });
  }
}
