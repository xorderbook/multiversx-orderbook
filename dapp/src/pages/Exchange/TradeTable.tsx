import { useEffect, useMemo, useState } from 'react';
import { useInterval } from '../../hooks/useInterval';
import { OpenOrderVo, TradeVo } from './Types';
import { useMarket } from 'store';
import { format } from '../../utils/format';
import { useQuery } from 'hooks/query/usePairQuery';
import { useGetAccountInfo } from 'hooks';
import { openOrders, myTrades } from './api';

import { useTable } from 'react-table';
import fakeData from './MOCK_DATA.json';
import { Address } from '@multiversx/sdk-core/out';
const TradeTable = () => {
  const { market } = useMarket();
  const [openOrderVo, setOpenOrderVo] = useState<any[]>([]);
  const [trades, setTrades] = useState<any[]>([]);
  const [tab, setTab] = useState('openOrders');
  // const { getAddressOrderIds } = useQuery();
  const { address } = useGetAccountInfo();

  const data = useMemo(() => fakeData, []);
  const columns = useMemo(
    () => [
      {
        Header: 'ID',
        accessor: 'id'
      },
      {
        Header: 'First Name',
        accessor: 'first_name'
      },
      {
        Header: 'Last Name',
        accessor: 'last_name'
      },
      {
        Header: 'Email',
        accessor: 'email'
      },
      {
        Header: 'Gender',
        accessor: 'gender'
      },
      {
        Header: 'University',
        accessor: 'university'
      }
    ],
    []
  );
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } =
    useTable({ columns, data });

  useEffect(() => {
    const fetchOrder = async () => {
      let openOrderList = await openOrders({
        marketID: market,
        address
      });
      setOpenOrderVo(openOrderList);
    };
    fetchOrder();
  }, []);

  useInterval(() => {
    const fetchOrder = async () => {
      let openOrderList = await openOrders({
        marketID: market,
        address
      });
      setOpenOrderVo(openOrderList);
    };

    const fetchTrade = async () => {
      let tradeList = await myTrades({
        marketID: market,
        address
      });
      setTrades(tradeList);
    };

    fetchOrder();
    fetchTrade();
  }, 2 * 1000);

  const convertToLocalTime = (timestamp: any) => {
    const date = new Date(timestamp * 1000);
    // Get the local time components from the Date object
    const localTime = {
      year: date.getFullYear(),
      month: date.getMonth() + 1, // Months are zero-based, so adding 1
      day: date.getDate(),
      hours: date.getHours(),
      minutes: date.getMinutes(),
      seconds: date.getSeconds()
    };

    // Format the local time components as a string
    const formattedTime = `${localTime.year}-${localTime.month}-${localTime.day} ${localTime.hours}:${localTime.minutes}:${localTime.seconds}`;
    return formattedTime;
  };

  const OpenOrderTable = () => {
    return (
      <div className='h-[17rem] overflow-y-auto text-gray-400 text-sm'>
        <div className='flex w-full items-center text-xs text-white py-2 px-3'>
          <p className='w-1/4 text-left md:w-[20%]'>Date</p>
          <p className='w-1/4 text-left md:w-[8%]'>Pair</p>
          <p className='w-1/4 text-center md:w-[8%]'>Type</p>
          <p className='w-1/4 text-center md:w-[8%]'>Side</p>
          <p className='hidden w-[12%] text-left md:flex'>Price</p>
          <p className='hidden w-[12%] text-left md:flex'>Amount</p>
          <p className='hidden w-[12%] text-left md:flex'>Filled</p>
          <p className='hidden w-[12%] text-left md:flex'>Total</p>
          <p className='hidden w-[12%] items-center justify-center text-center md:flex text-gray-text'>
            Cancel All
          </p>
        </div>
        {openOrderVo.map((order, index) => (
          <div
            key={index}
            className='p-3 flex w-full items-center text-xs text-white overflow-hidden'
          >
            <p className='w-1/4 text-left md:w-[20%]'>
              {convertToLocalTime(order.TimeStamp)}
            </p>
            <p className='w-1/4 text-left md:w-[8%]'>{market}</p>
            <p className='w-1/4 text-center md:w-[8%]'>{'limit'}</p>
            <p className='w-1/4 text-center md:w-[8%]'>
              {order.Side == 0 ? 'buy' : 'sell'}
            </p>
            <p className='hidden w-[12%] text-left md:flex'>
              {format(Number(order.Price))}
            </p>
            <p className='hidden w-[12%] text-left md:flex'>
              {format(Number(order.Qty - order.Filled))}
            </p>
            <p className='hidden w-[12%] text-left md:flex'>
              {format(Number(order.Filled))}
            </p>
            <p className='hidden w-[12%] text-left md:flex'>
              {format(Number(order.Qty))}
            </p>
            <p className='hidden w-[12%] items-center justify-center text-center md:flex text-gray-text'>
              <button className='bg-panel py-2 px-6 rounded-sm'>Cancel</button>
            </p>
          </div>
        ))}
      </div>
    );
  };

  const getBechAddress = (hexAddress: string) => {
    const address = Address.fromHex(hexAddress);
    return address.bech32();
  };

  const getSide = (order: any) => {
    if (getBechAddress(order.taker) == address) {
      return order.takerSide;
    }
    return order.takerSide == 'buy' ? 'sell' : 'buy';
  };

  const TradesTable = () => {
    return (
      <div className='h-[17rem] overflow-y-auto text-gray-400 text-sm'>
      <div className='flex w-full items-center text-xs text-white py-2 px-3'>
        <p className='w-1/4 text-left md:w-[20%]'>Date</p>
        <p className='w-1/4 text-left md:w-[8%]'>Pair</p>
        <p className='w-1/4 text-center md:w-[8%]'>Type</p>
        <p className='w-1/4 text-center md:w-[8%]'>Side</p>
        <p className='hidden w-[12%] text-left md:flex'>Price</p>
        <p className='hidden w-[12%] text-left md:flex'>Amount</p>
        <p className="hidden w-[12%] text-left md:flex">Executed</p>
      </div>
      {trades.map((order, index) => (
        <div
          key={index}
          className='p-3 flex w-full items-center text-xs text-white overflow-hidden'
        >
          <p className='w-1/4 text-left md:w-[20%]'>
            {convertToLocalTime(order.updatedAt)}
          </p>
          <p className='w-1/4 text-left md:w-[8%]'>{market}</p>
          <p className='w-1/4 text-center md:w-[8%]'>{'limit'}</p>
          <p className='w-[12%] text-center md:w-[8%]'>{getSide(order)}</p>
          <p className='hidden w-[12%] text-left md:flex'>
            {format(Number(order.price))}
          </p>
          <p className='hidden w-[12%] text-left md:flex'>
            {format(order.amount / 10 ** 18)}
          </p>
          <p className='hidden w-[12%] text-left md:flex'>
            <a href={order.hashUrl} target="_blank" className="hover:underline">Transaction</a>
          </p>
        </div>
      ))}
    </div>
  );
       
  };

  return (
    <div className=' '>
      <div className='flex items-center text-[#23F7DD] border-b border-[#27272B] bg-[#27272B] '>
        <div
          className={
            tab == 'openOrders'
              ? 'cursor-pointer px-3 py-2 text-center text-gray-text font-bold text-primary  border-[1px] underline  underline-offset-[14px] border-transparent'
              : 'cursor-pointer px-3 py-2 text-center text-gray-text text-primary'
          }
          onClick={() => {
            setTab('openOrders');
          }}
        >
          Open Orders ({openOrderVo.length})
        </div>
        <div
          className={
            tab == 'tradeHistory'
              ? 'cursor-pointer px-3 py-2 text-center text-gray-text font-bold border-[1px] underline  underline-offset-[14px] border-transparent'
              : 'cursor-pointer px-3 py-2 text-center text-gray-text ext-primary'
          }
          onClick={() => {
            setTab('tradeHistory');
          }}
        >
          Trade History ({trades.length})
        </div>
      </div>
      {tab == 'openOrders' ? OpenOrderTable() : TradesTable()}
    </div>
  );
};

export default TradeTable;
