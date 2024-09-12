import { Card } from 'components/Card';
import { contractAddress } from 'config';
import { AuthRedirectWrapper } from 'wrappers';
import TradePanel from './TradePanel';
import Orderbook from './Orderbook';
import TradeForm from './TradeForm';
import ExchangeHeader from 'components/Layout/Header/ExchangeHeader/ExchangeHeader';
import Footer from './Footer';
import { Header } from 'components/Layout/Header';

export const Exchange = () => (
  <div className='flex flex-col w-full bg-panel h-screen  '>
    <Header />

    <div className='flex flex-row h-full  '>
      <TradePanel></TradePanel>
      <Orderbook></Orderbook>
      <TradeForm></TradeForm>
    </div>
    <Footer></Footer>
  </div>
);
