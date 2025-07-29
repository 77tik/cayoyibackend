/*
 * @Author: harvey_li
 * @Date: 2022-08-09 09:23:55
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-09-02 13:15:19
 * @FilePath: \frontend\typings.d.ts
 * @Description:
 */
import '@umijs/max/typings';
import type { StyledInterface } from 'styled-components';

declare global {
  const styled: StyledInterface;
  interface Window {
    styled: StyledInterface;
  }
}
