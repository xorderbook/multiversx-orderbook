import React from 'react'

const Intro = () => {
    return (
        // <div className=' container mx-auto   '>
        //     <div className=' mx-auto py-[16rem]  w-full bg-[#171717]'>
        //         <div className="grid grid-cols-12 text-white ">

        //             <div className="col-span-12 lg:col-span-7 relative h-full flex justify-center">
        //                 <img className="shadow-green-200   rounded-2xl w-3/4 lg:w-full absolute scale-50 lg:-right-40 lg:top-1/2 lg:transform lg:-translate-y-1/2 xl:right-0 xl:w-[740px]" src="/book.png" alt="" />
        //             </div>

        //             <div className="col-span-12 lg:col-span-5 mb-12 lg:mb-0 relative z-10">
        //                 <h1 className="lg:text-lef text-6xl leading-[4.5rem] ">Full Order Book</h1>
        //                 <p className="text-center lg:text-left text-gray-400 text-2xl">
        //                     A complete limit order book exchange experience. Proceed here if you are used to trading in traditional finance or centralized exchanges.
        //                 </p>
        //             </div>
        //         </div>
        //     </div>
        // </div>
        <div className='container mx-auto  xl:max-w-[1280px] w-full'>
            < section id="features" className="flex md:flex-row flex-col sm:py-16 py-6" >
                <div className="flex-1 flex justify-center items-start flex-col">
                    <div className="col-span-12 lg:col-span-7 relative h-full flex justify-center">
                        <img className=" " src="/book.png" alt="" />
                    </div>

                    {/* <Button styles={`mt-10`} /> */}
                </div>



                <div className={`flex-1  flex justify-center items-star text-2xl font-semibold flex-col text-[#93ACB8] text-[18px]`}>
                    A complete limit order book exchange experience. Proceed here if you are used to trading in traditional finance or centralized exchanges.
                </div>
            </section >
        </div >


    )
}

export default Intro