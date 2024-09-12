const fs = require('fs');

const BTC_ESDT = "btcaddress"
const USDT_ESDT = "usdtAddress"
const ETH_ESDT = ""

const BTC_UDST_MARKET = ""
const ETH_USDT_MARKET = ""

function genEngineConfig() {
    const config = `
        [
            {
            "market": "BTC-USDT",
            "exchangeAddress": "${BTC_UDST_MARKET}",
            "baseSymbol": "BTC",
            "quoteSymbol": "USDT",
            "baseESDT": "${BTC_ESDT}",
            "quoteESDT": "${USDT_ESDT}"
            },
            {
            "market": "ETH-USDT",
            "exchangeAddress": "${ETH_USDT_MARKET}",
            "baseSymbol": "ETH",
            "quoteSymbol": "USDT",
            "baseESDT": "${ETH_ESDT}",
            "quoteESDT": "${USDT_ESDT}"
            }
      ]
    `
    fs.writeFile('engineConfig.txt', config, (err) => {
        if (err) {
            console.error('Error writing file:', err);
        } else {
            console.log('File written successfully.');
        }
    });
}

function genDappConfig() {
    const config = `{
        market: [
          {
            market: "BTC-USDT",
            exchangeAddress:"${BTC_UDST_MARKET}",
            baseSymbol: "BTC",
            quoteSymbol: "USDT",
            baseESDT: "${BTC_ESDT}”,
            quoteESDT: "${USDT_ESDT}",
            logo: bitcoin
          },
          {
            market: "ETH-USDT",
            exchangeAddress:"${ETH_USDT_MARKET}",
            baseSymbol: "ETH",
            quoteSymbol: "USDT",
            baseESDT: "${ETH_ESDT}",
            quoteESDT: "${USDT_ESDT}",
            logo: ethereum
          }
        ],
        ESDTList: [
          {
            Symbol: "USDT",
            ESDT: "${USDT_ESDT}",
            logo: usdt,
          },
          {
            Symbol: "BTC",
            ESDT: "${BTC_ESDT}",
            logo: bitcoin
          },
          {
            Symbol: "ETH",
            ESDT: "${ETH_ESDT}",
            logo: ethereum
          },
        ],
        faucetAddress: "erd1qqqqqqqqqqqqqpgqe47azwqu97maz0f2eu86qqwfvymlwgvsymtsw0veu2",
      }`;

    fs.writeFile('dappConfig.txt', config, (err) => {
        if (err) {
            console.error('Error writing file:', err);
        } else {
            console.log('File written successfully.');
        }
    });

}

genDappConfig()
genEngineConfig()