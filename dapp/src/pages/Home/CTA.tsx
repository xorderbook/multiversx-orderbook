import styles from './style';
// import Button from "./Button";

const CTA = () => (
  <div className={`bg-primary ${styles.paddingX} ${styles.flexCenter} py-14`}>
    <div className={`${styles.boxWidth}`}>
      <div className={`flex justify-center items-center flex-col ${styles.marginY} ${styles.padding} sm:flex-row  bg-black-gradient-2 rounded-[20px] box-shadow`}>
        <div className='flex-1 flex flex-col justify-center'>
          <h2 className='font-poppins font-semibold xs:text-[48px] text-[40px] text-white xs:leading-[76.8px] leading-[66.8px] w-full text-center'>
            Let’s start trading now!
          </h2>
          <div className='mx-auto mt-8'>
            <button className='px-4 py-2 bg-[#21F7DC] rounded-sm'>
              Trade now
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
);

export default CTA;
