import React, { ChangeEvent, useRef } from 'react';
import { useEffect, useState } from 'react';
import { useSendPairTransaction } from 'hooks/transactions/useSendPairTransaction';
import { SessionEnum } from 'localConstants';
import { useMarket } from 'store';
import { useInterval } from 'hooks/useInterval';
import axios from 'axios';
import { contractAddress } from 'config/config.contract';
import { format } from '../../utils/format';
import { useGetAccountInfo } from 'hooks';
import Fee from './Fee';
import AmountSlider from './AmountSlider';
import BTC from '../../assets/img/BTC.svg';
import usdt from '../../assets/img/usdt.svg';
import ETH from '../../assets/img/ETH.svg';

const TradeForm = () => {
  const [price, setPrice] = useState('');
  const [amount, setAmount] = useState('');
  const [side, setSide] = useState(0);
  const [quote, setQuote] = useState('');
  const [base, setBase] = useState('');
  const [quoteBalance, setQuoteBalance] = useState('');
  const [baseBalance, setBaseBalance] = useState('');
  const { market, sliderPercent, setBaseWalletBalance, setQuoteWalletBalance } = useMarket();
  const { address } = useGetAccountInfo();

  const [baseLogo, setBaseLogo] = useState<any>();
  const [quoteLogo, setQuoteLogo] = useState<any>();

  const { sendCreateBuyOrder, sendCreateSellOrder, transactionStatus } =
    useSendPairTransaction(SessionEnum.abiPingPongSessionId);

  const getMarketCfg = (market: string) => {
    return contractAddress.market.filter((item) => item.market === market)[0];
  };

  const placeOrder = async () => {
    if (side == 0) {
      await sendCreateBuyOrder(market, Number(amount), Number(price));
    } else {
      await sendCreateSellOrder(market, Number(amount), Number(price));
    }
  };

  const changeSide = (side: number) => {
    setSide(side);
  };

  const fetchBalance = async () => {
    axios
      .get(`https://devnet-gateway.multiversx.com/address/${address}/esdt`)
      .then((response) => {
        // Handle the successful response
        let esdts = response.data.data.esdts;
        setBaseBalance(
          format(
            esdts[getMarketCfg(market).baseESDT]
              ? esdts[getMarketCfg(market).baseESDT].balance / 10 ** 18
              : 0
          )
        );
        setQuoteBalance(
          format(
            esdts[getMarketCfg(market).quoteESDT]
              ? esdts[getMarketCfg(market).quoteESDT].balance / 10 ** 18
              : 0
          )
        );

        let walletBaseBalance = format(
          esdts[getMarketCfg(market).baseESDT]
            ? esdts[getMarketCfg(market).baseESDT].balance / 10 ** 18
            : 0
        );
        let walletQuoteBalance = format(
          esdts[getMarketCfg(market).quoteESDT]
            ? esdts[getMarketCfg(market).quoteESDT].balance / 10 ** 18
            : 0
        );

        // setBaseWalletBalance(Number(walletBaseBalance));
        //setQuoteWalletBalance(Number(walletQuoteBalance));
      })
      .catch((error) => {
        // Handle the error
        console.error(error);
      });
  };

  useInterval(() => {
    fetchBalance();
  }, 1 * 1000);

  useEffect(() => {
    let arr = market.split('-');
    const base = arr[0];
    const quote = arr[1];
    setBase(base);
    setQuote(quote);
    fetchBalance();

    if (base == "BTC") {
      setBaseLogo(BTC)
    } else if (base == "ETH") {
      setBaseLogo(ETH)
    }

    setQuoteLogo(usdt)
  }, [market]);

  useEffect(() => {
    if (sliderPercent == 0) {
      setAmount('')
      return
    }
    let baseActualVal = (sliderPercent * Number(baseBalance) / 100)
    setAmount(baseActualVal.toString());
  }, [sliderPercent]);

  function handleAmountChange(event: ChangeEvent<HTMLInputElement>) {
    const currentVal = amount;
    const val = event.target.value;
    const inputValue = val.replace(/-/g, '');
    const regex = /^\d*\.?\d*$/;
    if (regex.test(inputValue)) {
      setAmount(inputValue);
    } else {
      setAmount(currentVal);
    }
  }

  function handlePriceChange(event: ChangeEvent<HTMLInputElement>) {
    const currentVal = price;
    const val = event.target.value;
    const inputValue = val.replace(/-/g, '');
    const regex = /^\d*\.?\d*$/;
    if (regex.test(inputValue)) {
      setPrice(inputValue);
    } else {
      setPrice(currentVal);
    }
  }

  const getTokenLogo = () => {
    return BTC

  }

  return (
    <div className='p-3 flex  w-[20%] flex-col gap-y-4 overflow-y-hidden   border-l border-gray-btn bg-gray-bg  bg-main border-[#3E3D40] text-white  '>
      Balances
      <div>
        <div className='flex justify-between text-sm'>
          <span>Asset</span>
          <span>Wallet</span>
        </div>

        <div className='mt-2 flex justify-between text-xs  text-gray-400'>
          <span className='flex space-x-2 justify-center items-center'>
            <img src={baseLogo} alt='icon' className='w-[21px]' />
            <span>{base}</span>
          </span>
          <span>{baseBalance}</span>
        </div>

        <div className='flex justify-between text-xs text-gray-400 mt-1'>
          <div className='flex space-x-2 justify-center items-center'>
            <img src={usdt} alt='icon' className='w-[21px]' />
            <span>{quote}</span>
          </div>
          <span>{quoteBalance}</span>
        </div>
      </div>
      <div className='flex justify-between p-2 w-full text-black'>
        <button
          onClick={() => {
            changeSide(0); // buy
          }}
          className={
            side == 0
              ? 'bg-[#23F7DD] w-full py-2 rounded-s'
              : 'bg-[#27272B] w-full py-2 rounded-s text-white'
          }
        >
          Buy
        </button>
        <button
          onClick={() => {
            changeSide(1); // sell
          }}
          className={
            side == 1
              ? 'bg-red-400 w-full py-2 rounded-e'
              : 'bg-[#27272B] w-full py-2 rounded-e text-white'
          }
        >
          Sell
        </button>
      </div>
      <a className='font-semibold text-primary'>Limit</a>
      <div className='flex w-full flex-col rounded border bg-[#27272B] px-2 py-3 border-transparent'>
        <div className='mt-2 flex items-center text-xs'>
          <p className='text-gray-500'>Price</p>
          <input
            pattern='[0-9]*'
            placeholder='0.0'
            type='text'
            className='mx-1 w-full flex-1 bg-transparent text-right outline-none [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none'
            value={price}
            onChange={handlePriceChange}
          />
          <p>{quote}</p>
        </div>
      </div>
      <div className='flex w-full flex-col rounded border bg-[#27272B] px-2 py-3 border-transparent'>
        <div className='mt-2 flex items-center text-xs'>
          <p className='text-gray-500'>Amount</p>
          <input
            pattern='[0-9]*'
            placeholder='0.0'
            type='text'
            className='mx-1 w-full flex-1 bg-transparent text-right outline-none [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none'
            value={amount}
            onChange={handleAmountChange}
          />
          <p>{base}</p>
        </div>
      </div>
      <AmountSlider></AmountSlider>
      <Fee></Fee>
      <div className='w-full mt-10 px-3'>
        <button
          className={
            side == 0
              ? 'text-black bg-[#23F7DD] py-2 w-full rounded'
              : 'text-black bg-red-400 py-2 w-full rounded'
          }
          onClick={placeOrder}
        >
          {side == 0 ? 'Buy' : 'Sell'}
        </button>
      </div>
    </div>
  );
};

export default TradeForm;
