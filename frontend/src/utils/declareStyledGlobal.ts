/*
 * @Author: harvey_li
 * @Date: 2022-09-02 11:58:40
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-09-02 13:18:32
 * @FilePath: \frontend\src\utils\declareStyledGlobal.ts
 * @Description: 在此处引入styled 并赋值给window.styled
 */

import styled from 'styled-components';
(function (window, styled) {
  window.styled = styled;
})(window, styled);
