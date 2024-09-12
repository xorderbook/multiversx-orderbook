import ethereum from "../assets/img/ethereum.svg";
import usdc from "../assets/img/usdc.svg";
import bitcoin from "../assets/img/bitcoin.svg"
import usdt from "../assets/img/usdt.svg";

export const contractAddress = {
  market: [
    {
      market: 'BTC-USDT',
      exchangeAddress:
        'erd1qqqqqqqqqqqqqpgqx23lca0gfkvnnqclfseasq0pdctfxst4ymtsjsedkm',
      baseSymbol: 'BTC',
      quoteSymbol: 'USDT',
      baseESDT: 'BTC-187d25',
      quoteESDT: 'USDT-cf0380',
      logo: bitcoin
    },
    {
      market: 'ETH-USDT',
      exchangeAddress:
        'erd1qqqqqqqqqqqqqpgqp9szauyjdfdx8qgwn22h3sv4p3tznz6vymtsgq9s09',
      baseSymbol: 'ETH',
      quoteSymbol: 'USDT',
      baseESDT: 'ETH-dd05ee',
      quoteESDT: 'USDT-cf0380',
      logo: ethereum
    }
  ],
  ESDTList: [
    {
      Symbol: "USDT",
      ESDT: 'USDT-cf0380',
      logo: usdt,
    },
    {
      Symbol: "BTC",
      ESDT: 'BTC-187d25',
      logo: bitcoin
    },
    {
      Symbol: "ETH",
      ESDT: 'ETH-dd05ee',
      logo: ethereum
    },
  ],
  faucetAddress: "erd1qqqqqqqqqqqqqpgqram9w26rgsr2r4rerll5xuym9lzlrctkymts3ntfpk",
};
