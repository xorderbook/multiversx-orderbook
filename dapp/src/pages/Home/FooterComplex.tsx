import React from 'react'
import { people01, people02, people03, facebook, instagram, linkedin, twitter, airbnb, binance, coinbase, dropbox, send, shield, star } from "../../assets/demo"
import logo from '../../assets/img/xorderbook.png';
const footerLinks = [
    {
        title: "Useful Links",
        links: [
            {
                name: "Content",
                link: "https://www.hoobank.com/content/",
            },
            {
                name: "How it Works",
                link: "https://www.hoobank.com/how-it-works/",
            },
            {
                name: "Create",
                link: "https://www.hoobank.com/create/",
            },
            {
                name: "Explore",
                link: "https://www.hoobank.com/explore/",
            },
            {
                name: "Terms & Services",
                link: "https://www.hoobank.com/terms-and-services/",
            },
        ],
    },
    {
        title: "Community",
        links: [
            {
                name: "Twitter",
                link: " ",
            },
            {
                name: "Discord",
                link: " ",
            },
            {
                name: "Telegram",
                link: "  ",
            },
            {
                name: "Media",
                link: " ",
            },
           
        ],
    },
    {
        title: "Product",
        links: [
            {
                name: "Exchange",
                link: " ",
            },
            {
                name: "Farm",
                link: " ",
            },
        ],
    },
];


const socialMedia = [
    {
        id: "social-media-3",
        icon: twitter,
        link: "https://www.twitter.com/",
    },
];

const Footer = () => {
    return (
        <div className='container mx-auto  xl:max-w-[1280px] w-full'>
            <section className={`flex justify-center items-center sm:py-16 py-6 flex-col`}>
                <div className={`flex justify-center items-start md:flex-row flex-col mb-8 w-full`}>
                    <div className="flex-[1] flex flex-col justify-start mr-10">
                        <div className='flex flex-row items-center'>
                            <img
                                src={logo}
                                alt="hoobank"
                                className="  h-[72.14px] object-contain  "

                            />
                            <span className='ml-2 text-white text-3xl '> xOrderbook  </span>
                        </div>
                        <p className={`font-poppins font-normal  text-white text-[18px] leading-[30.8px] mt-4 max-w-[312px]`}>
                            A orderbook DEX built on MultiversX.
                        </p>
                    </div>

                    <div className="flex-[1.5] w-full flex flex-row justify-between flex-wrap md:mt-0 mt-10">
                        {footerLinks.map((footerlink) => (
                            <div key={footerlink.title} className={`flex flex-col ss:my-0   min-w-[150px]`}>
                                <h4 className="font-poppins font-medium text-[18px] leading-[27px] text-white">
                                    {footerlink.title}
                                </h4>
                                <ul className="list-none mt-4">
                                    {footerlink.links.map((link, index) => (
                                        <li
                                            key={link.name}
                                            className={`font-poppins font-normal text-[16px] leading-[24px] text-[#B3B4B7] hover:text-secondary cursor-pointer ${index !== footerlink.links.length - 1 ? "mb-4" : "mb-0"
                                                }`}
                                        >
                                            {link.name}
                                        </li>
                                    ))}
                                </ul>
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