// import { Header } from 'components/Layout/Header';
import { AuthRedirectWrapper, PageWrapper } from 'wrappers';
import { Header } from './Header';
import Intro from './Intro';
//import Footer from './Footer';
import Footer from './FooterComplex';
import { Advantage } from './Advantage';
import GetStarted from './GetStarted';
import CTA from './CTA';
import Billing from './Billing';
import CryptoPrice from './CryptoPrice';

export const Home = () => {
  return (
    <div className='w-full bg-[#171717]'>
      <Header></Header>
      <div className='bg-[#171717]  xl:max-w-[1280px] w-full mx-auto  '>

        <div className=' mx-auto py-[8rem]  bg-[#171717]'>
          <div className="grid grid-cols-12 text-white ">
            <div className="col-span-12 lg:col-span-5 mb-12 lg:mb-0 relative z-10">
              <h1 className="mb-6 lg:text-lef   leading-[6.5rem]  font-semibold ss:text-[72px] text-[68px] text-gradient">
                <div className="  md:mr-4 mr-0 flex flex-row">
                  Orderbook
                  {/* <GetStarted /> */}
                </div>

                Trading on
                Multivers
                <span className='text-[#23F7DD]'>X</span>
              </h1>
              <p className="intro-p mb-10 text-center lg:text-left text-gray-400 text-xl">
                Experience the
                blazing fast
                trading speed at<br />
                xOrderbook.
              </p>
            </div>
            <div className="col-span-12 lg:col-span-7   h-full flex justify-center relative">
              <img className="  shadow-zinc-800 shadow   w-3/4 lg:w-full absolute   lg:-right-40 lg:top-1/2 lg:transform lg:-translate-y-1/2 xl:right-0 xl:w-[680px] h-[460px] " src="/bg.png" alt="" />

              <div className="absolute z-[0] w-[40%] h-[35%] top-0 pink__gradient" />
              {/* <div className="absolute z-[1] w-[80%] h-[80%] rounded-full white__gradient bottom-40" /> */}
              {/* <div className="absolute z-[0] w-[50%] h-[50%] right-20 bottom-20 blue__gradient" /> */}

            </div>
          </div>
        </div>
      </div>
      {/* <Intro></Intro> */}

      <CryptoPrice />
      <Advantage></Advantage>
      <Billing></Billing>
      <CTA />
      <Footer></Footer>
    </div>
  );
};
