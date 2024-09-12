import React, { useState } from 'react';
import { useMarket } from 'store';

const AmountSlider = () => {
  const [value, setValue] = useState(0);

  const {  setSliderPercent } = useMarket();

  const handleRangeChange = (event: any) => {
    const percent = event.target.value;
    console.log('', event.target.value);
    setValue(percent);
    setSliderPercent(percent)
  };

  return (
    <div className='w-full'>
      <input
        type='range'
        min={0}
        max='100'
        className='range w-full  accent-[#27272C]'
        step='25'
        value={value}
        onChange={handleRangeChange}
      />

      {/* <input type='range'min={0} max='100' value='75' className='range w-full  accent-[#27272C]' step='25'/> */}

      {/* <input type="range" className="h-2 w-full cursor-ew-resize appearance-none rounded-full bg-gray-200 disabled:cursor-not-allowed"> */}
      <div className='w-full flex justify-between text-xs px-2'>
        <span>0</span>
        <span>25</span>
        <span>50</span>
        <span>75</span>
        <span>100</span>
      </div>
    </div>
  );
};

export default AmountSlider;
