import TradingView from './TradingView';
import TradeTable from './TradeTable';
import ExchangeHeader from 'components/Layout/Header/ExchangeHeader/ExchangeHeader';

const TradePanel = () => {
  return (
    <div className='flex h-full w-[60%] flex-col  overflow-y-hidden border-gray-btn bg-gray-bg  border-r border-[#3E3D40]'>
      <div className='flex items-center justify-between bg-dark px-4 md:px-6 md:py-3 bg-main border-b border-[#3E3D40] w-full'>
        <div className='flex items-center text-gray-300 w-full'><ExchangeHeader></ExchangeHeader></div>
      </div>
      <div className='h-[60%] bg-main'><TradingView></TradingView></div>
      <div className='flex-1 bg-main z-30 '>
        <TradeTable></TradeTable>
      </div>
    </div>
  );
};

export default TradePanel;
