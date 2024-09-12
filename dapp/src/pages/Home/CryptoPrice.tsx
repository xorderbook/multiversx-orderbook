import React, { useState, useEffect } from 'react';
import BTC from '../../assets/img/BTC.svg';
import ETH from '../../assets/img/ETH.svg';
import EGLD from '../../assets/img/EGLD.svg';
import LTC from '../../assets/img/LTC.svg';
import SOL from '../../assets/img/SOL.svg';
import Solana from '../../assets/img/Solana.svg';
import { useInterval } from 'hooks/useInterval';

const CryptoPrice = () => {
  const [price, setPrice] = useState<number[]>([]);
  const [lastPrice, setLastPrice] = useState<number[]>([]);

  const coins = [
    {
      id: 'btc',
      name: 'Bitcoin',
      icon: BTC,
      link: 'https://www.instagram.com/',
      price: 0,
      cls: ""
    },
    {
      id: 'eth',
      name: 'Ethereum',
      icon: ETH,
      link: 'https://www.facebook.com/',
      price: 0,
      cls: ""
    },
    {
      id: 'sol',
      name: 'Solana',
      icon: Solana,
      link: 'https://www.linkedin.com/',
      price: 0,
      cls: ""
    },
    {
      id: 'egld',
      name: 'MultiversX',
      icon: EGLD,
      link: 'https://www.twitter.com/',
      price: 0,
      cls: ""
    },
  ];

  const getPrice = () => {
    fetch(`https://api.binance.com/api/v3/ticker/price?symbols=["BTCUSDT","ETHUSDT","EGLDUSDT","SOLUSDT"]`)
      .then((res) => res.json())
      .then((data) => {
        setLastPrice(price)

        const prices = data.map((item: { price: any; }) => Number(item.price).toFixed(2));
        //  [prices[2], prices[3]] = [prices[3], prices[2]];
        setPrice(prices)
      })
      .catch((err) => console.log(err));
  };

  useEffect(() => {
    getPrice();
  }, []);

  const updatedCoins = coins.map((coin, index) => {
    return {
      ...coin,
      price: price[index],
      cls: price[index] >= lastPrice[index] ? "text-green-400" : "text-[#EB524F]"
    };
  });

  useInterval(() => {
    getPrice()
  }, 2 * 1000);

  return (
    <div className='mx-auto xl:max-w-[1280px] w-full flex flex-row py-20'>
      <div className='w-full grid grid-cols-4'>
        {updatedCoins.map((coin, index) => (
          <div
            key={index}
            className='flex flex-row items-center text-white space-x-6'
          >
            <img src={coin.icon} alt='icon' className='w-[42px]' />
            <span className={coin.cls}>{coin.name} {coin.price}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CryptoPrice;
