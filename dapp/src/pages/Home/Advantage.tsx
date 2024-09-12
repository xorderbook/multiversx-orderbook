import React from 'react'
import { people01, people02, people03, facebook, instagram, linkedin, twitter, airbnb, binance, coinbase, dropbox, send, shield, star } from "../../assets/demo"
const features = [
  {
    id: "feature-1",
    icon: star,
    title: "Fast",
    content:
      "High transactional throughput with a 1-2 second finality.",
  },
  {
    id: "feature-2",
    icon: shield,
    title: "Trade trustlessly",
    content:
      "Completely trustless and self-custodial.",
  },
  {
    id: "feature-3",
    icon: send,
    title: "Low cost",
    content:
      "Low fee rates allowing users to maximize their returns.",
  },
];



const FeatureCard = ({ icon, title, content, index }: any) => (
  <div className={`flex flex-row p-6 rounded-[20px] ${index !== features.length - 1 ? "mb-6" : "mb-0"} feature-card`}>
    <div className={`w-[64px] h-[64px] rounded-full flex justify-center items-center bg-[#0A1922]`}>
      <img src={icon} alt="star" className="w-[50%] h-[50%] object-contain" />
    </div>
    <div className="flex-1 flex flex-col ml-3">
      <h4 className="font-poppins font-semibold text-white text-[18px] leading-[23.4px] mb-1">
        {title}
      </h4>
      <p className="font-poppins font-normal text-[#B5B4B9] text-[16px] leading-[24px]">
        {content}
      </p>
    </div>
  </div>
);

export const Advantage = () => {
  return (
    <div className='  mx-auto  xl:max-w-[1280px] w-full'>
      < section id="features" className="flex md:flex-row flex-col sm:py-16 py-6" >
        <div className="flex-1 flex justify-center items-start flex-col h-full">
          <h2 className="font-poppins font-semibold xs:text-[48px] text-[40px] text-white xs:leading-[76.8px] leading-[66.8px] w-full">
            Next Generation of Exchange <br className="sm:block hidden" />
          </h2>
          <p className={` font-poppins font-normal text-[#93ACB8] text-[18px] leading-[30.8px] max-w-[470px] mt-5`}>
            Combines the feel and easy experience of a CEX with the self-custody of a DEX.
          </p>

          {/* <Button styles={`mt-10`} /> */}
        </div>

        <div className={`flex-1 flex justify-center items-start md:ml-10 ml-0 md:mt-0 mt-10 relative flex-col`}>
          {features.map((feature, index) => (
            <FeatureCard key={feature.id} {...feature} index={index} />
          ))}
        </div>
      </section >
    </div >
  )
}
