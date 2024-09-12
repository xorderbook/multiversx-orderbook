import React, { useEffect, useState } from 'react';
import { SessionEnum } from 'localConstants';
import { useMintTransaction } from 'hooks/transactions/useSendFaucetTransaction';
import { contractAddress } from '../../config/config.contract';
import { useGetAccountInfo } from 'hooks';
import axios from 'axios';
import { format } from '../../utils/format';
import { Header } from 'components/Layout/Header';
import { useInterval } from 'hooks/useInterval';

const ESDT_LIST = contractAddress.ESDTList;
const Faucet = () => {
  const { address } = useGetAccountInfo();
  const [balance, setBalance] = useState<any>([]);

  const { sendMintTransaction, mintTokenTransactionStatus } = useMintTransaction(SessionEnum.abiPingPongSessionId);
  
  const onMint = async (esdtId: string, symbol: string) => {
    if (symbol == 'BTC') {
      await sendMintTransaction(esdtId, '10000000000000000000');
    } else if (symbol == 'USDT') {
      await sendMintTransaction(esdtId, '10000000000000000000000');
    } else if (symbol == 'ETH') {
      await sendMintTransaction(esdtId, '100000000000000000000');
    }
  };

  const getAssetLogo = (symbol: string) => {
    return ESDT_LIST.filter((item) => item.Symbol === symbol)[0].logo;
  };

  const fetchBalance = async () => {
    let faucetBalance: any = [];
    let myBalance: any = [];
    axios
      .get(`https://devnet-gateway.multiversx.com/address/${address}/esdt`)
      .then((response) => {
        let esdts = response.data.data.esdts;
        // console.log('account esdts', esdts);
        ESDT_LIST.map((token) => {
          let ret = Object.keys(esdts).filter((i: any) => i === token.ESDT);
          if (ret[0]) {
            const merged = Object.assign(esdts[ret[0]], token);
            myBalance.push(merged);
          }
        });

        axios
          .get(
            `https://devnet-gateway.multiversx.com/address/${contractAddress.faucetAddress}/esdt`
          )
          .then((response) => {
            let esdts = response.data.data.esdts;
            // console.log('faucet ', contractAddress.faucetAddress);
            // console.log('faucet esdts', esdts);
            ESDT_LIST.map((token) => {
              let ret = Object.keys(esdts).filter((i: any) => i === token.ESDT);
              const merged = Object.assign(esdts[ret[0]], token);
              faucetBalance.push(merged);
            });

            faucetBalance.forEach((item: any) => {
              myBalance.forEach((y: any) => {
                if (item.ESDT == y.ESDT) {
                  item['my'] = y;
                }
              });
            });
            setBalance(faucetBalance);
          })
          .catch((error) => {
            console.error(error);
          });
      })
      .catch((error) => {
        console.error(error);
      });
  };

  useEffect(() => {
    fetchBalance();
  }, []);

  useInterval(() => {
    fetchBalance();
  }, 1 * 1000);

  return (
    <div className='w-full'>
      <Header />

      <div className='w-full bg-[#27272B] pt-4 h-full'>
        <div className='px-20  '>
          <table className='table-auto w-[1280px] mx-auto justify-start text-white  rounded-xl bg-[#000000]'>
            <thead className='h-16'>
              <tr className='text-left py-4'>
                <th>
                  <span className='ml-8'>Asset</span>
                </th>
                <th className=''>Mint Account</th>
                <th className=''>My balance </th>
              </tr>
            </thead>
            <tbody className='mt-20 rounded-xl'>
              {balance.map((item: any) => {
                return (
                  <tr
                    key={item.Symbol}
                    className='bg-[#202020] h-20 rounded-xl'
                  >
                    <td>
                      <div className='flex flex-row items-center'>
                        <img
                          className='h-8 ml-6'
                          src={getAssetLogo(item.Symbol)}
                          alt=''
                        />
                        <span className='mx-4'>{item.Symbol}</span>
                      </div>
                    </td>
                    <td>{format(item.balance / 10 ** 18)}</td>
                    <td>
                      <span>
                        {item.my ? format(item.my.balance / 10 ** 18) : 0}
                      </span>
                    </td>
                    <td className='mx-8'>
                      <div className='flex flex-row justify-center'>
                        <button
                          className='bg-[#21F7DC] px-16 py-2 rounded text-black'
                          onClick={() => {
                            onMint(item.ESDT, item.Symbol);
                          }}
                        >
                          Claim
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Faucet;
