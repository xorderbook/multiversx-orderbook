import { create } from 'zustand'
import { contractAddress } from '../config/config.contract';
const marketList = contractAddress.market

export interface Market {
    market: string,
    baseWalletBalance: number;
    quoteWalletBalance: number,
    sliderPercent: number,
    setMarket: (id: string) => void;
    setBaseWalletBalance: (balance: number) => void,
    setQuoteWalletBalance: (balance: number) => void,
    setSliderPercent: (balance: number) => void,
}

export const useMarket = create<Market>((set) => ({
    market: marketList[0].market,
    baseWalletBalance: 0,
    quoteWalletBalance: 0,
    sliderPercent: 0,
    setMarket: (market: string) => set(() => ({ market: market })),
    setBaseWalletBalance: (balance: number) => set(() => ({ baseWalletBalance: balance })),
    setQuoteWalletBalance: (balance: number) => set(() => ({ quoteWalletBalance: balance })),
    setSliderPercent: (balance: number) => set(() => ({ sliderPercent: balance })),
}))