import React from 'react'
import { facebook, instagram, linkedin, twitter } from "../../assets/demo"
import logo from '../../assets/img/xorderbook.png';
const footerLinks = [
    {
        title: "App",
    },
    {
        title: "Markets",
    },
    {
        title: "Docs",
    },
];

const socialMedia = [
    {
        id: "social-media-1",
        icon: instagram,
        link: "https://www.instagram.com/",
    },
    {
        id: "social-media-2",
        icon: facebook,
        link: "https://www.facebook.com/",
    },
    {
        id: "social-media-3",
        icon: twitter,
        link: "https://www.twitter.com/",
    },
    {
        id: "social-media-4",
        icon: linkedin,
        link: "https://www.linkedin.com/",
    },
];

const Footer = () => {
    return (
        <div className='container mx-auto xl:max-w-[1280px] w-full'>
            <section className={`flex justify-center items-center sm:py-16 py-6 flex-col`}>
                <div className={`flex justify-center items-start md:flex-row flex-col mb-8 w-full`}>
                    <div className="flex-[2] flex flex-col justify-start mr-10">
                        <div className='flex flex-row items-center'>
                            <img src={logo} alt="hoobank" className="h-[2.35rem] object-contain" />
                            <span className='ml-2 text-white text-2xl font-bold'> xOrderbook </span>
                        </div>
                    </div>
                    <div className="flex-[1.5] w-full flex flex-row justify-between flex-wrap md:mt-0 mt-10">
                        {footerLinks.map((footerlink) => (
                            <div key={footerlink.title} className={`flex flex-col ss:my-0 min-w-[80px]`}>
                                <h4 className="font-poppins font-medium text-[18px] leading-[27px] text-white">
                                    {footerlink.title}
                                </h4>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="w-full flex justify-between items-center md:flex-row flex-col pt-6 border-t-[1px] border-t-[#3F3E45]">
                    <p className="font-poppins font-normal text-center text-[18px] leading-[27px] text-white">
                        Copyright Ⓒ 2023 MultiversX Orderbook. All Rights Reserved.
                    </p>

                    <div className="flex flex-row md:mt-0 mt-6">
                        {socialMedia.map((social, index) => (
                            <img
                                key={social.id}
                                src={social.icon}
                                alt={social.id}
                                className={`w-[21px] h-[21px] object-contain cursor-pointer ${index !== socialMedia.length - 1 ? "mr-6" : "mr-0"}`}
                                onClick={() => window.open(social.link)}
                            />
                        ))}
                    </div>
                </div>
            </section>
        </div>
    )
}

export default Footer