/*
 * @Author: Even
 * @Date: 2022-03-28 14:13:26
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-08-05 17:14:40
 * @FilePath: \frontend\src\store\index.ts
 */

import { makeAutoObservable } from 'mobx';

class Store {
  constructor() {
    makeAutoObservable(this);
  }
  fontSize = 1;
  updateFontSize(fontSize: number) {
    this.fontSize = fontSize;
  }
  widthRate = 1
   updateWidthRate(widthRate: number) {
    this.widthRate = widthRate;
  }
  socketData: {
    page_id: number;
    date: number;
  } = {
    page_id: 0,
    date: 0,
  };
  updateSocketData(data: { page_id: number; date: number; }) {
    this.socketData = data;
  }
}

export default new Store();
