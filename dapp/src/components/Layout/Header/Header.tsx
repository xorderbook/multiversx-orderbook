import { Button } from 'components/Button';
import { MxLink } from 'components/MxLink';
import { environment } from 'config';
import { logout } from 'helpers';
import { useGetIsLoggedIn } from 'hooks';
import { RouteNamesEnum } from 'localConstants';
import LogoImg from '../../../assets/img/xorderbook.png';

export const Header = () => {
  const isLoggedIn = useGetIsLoggedIn();

  const handleLogout = () => {
    sessionStorage.clear();
    logout(`${window.location.origin}/unlock`, undefined, false);
  };

  return (
    <header className='flex flex-row align-center justify-between pl-6 pr-6 py-2 bg-[#171717] items-center'>
      <MxLink
        className='flex items-center justify-between'
        to={RouteNamesEnum.home}
      >
        <img className='h-[2.35rem] pr-3' src={LogoImg} alt='logo' />
        <span className='text-white text-2xl font-bold'>
          xOrderbook
        </span>
      </MxLink>

      <nav className='h-full w-full text-sm sm:relative sm:left-auto sm:top-auto sm:flex sm:w-auto sm:flex-row sm:justify-end sm:bg-transparent'>
        <div className='flex justify-end container mx-auto items-center gap-2'>
          <MxLink
            to={RouteNamesEnum.exchange}
            className='inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0   text-white ml-2  mr-3'
          >
            Trade
          </MxLink>

          <MxLink
            to={RouteNamesEnum.farm}
            className='inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0   text-white ml-2  mr-5'
          >
            Farm
          </MxLink>

          <MxLink
            to={RouteNamesEnum.faucet}
            className='inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0   text-white ml-2  mr-5'
          >
            Faucet
          </MxLink>
          <div className="border-r h-[60%] border-[#383838]"></div>
          <div className='flex gap-1 items-center'>
            <div className='w-2 h-2 rounded-full bg-green-500' />
            <p className='text-gray-600'>{environment}</p>
          </div>
          {isLoggedIn ? (
            <Button
              onClick={handleLogout}
              className='inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0 text-gray-600 hover:bg-slate-100 mx-0'
            >
              Close
            </Button>
          ) : (
            <MxLink to={RouteNamesEnum.unlock}>Connect</MxLink>
          )}
        </div>
      </nav>
    </header>
  );
};
