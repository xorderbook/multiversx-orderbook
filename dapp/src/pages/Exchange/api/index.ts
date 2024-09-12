import axios from "axios";
import { OpenOrderResp, OrderbookResp, TradeResp } from "../Types";
// import { BaseUrl } from "config";


export const BaseUrl = "http://127.0.0.1:8000/api/"; // local
export async function orderbook(data: any) {
  let resp = await axios.request<OrderbookResp>({
    url: `${BaseUrl}orderbook`,
    method: "post",
    data,
  })

  if (resp.data.data) {
    return JSON.parse(resp.data.data)
  }
  return null
}

export async function openOrders(data: any) {
  let resp = await axios.request<OpenOrderResp>({
    url: `${BaseUrl}openOrder`,
    method: "post",
    data,
  })
  return resp.data.data
}

export async function myTrades(data: any) {
  let resp = await axios.request<TradeResp>({
    url: `${BaseUrl}trades`,
    method: "post",
    data,
  })
  return resp.data.data
}