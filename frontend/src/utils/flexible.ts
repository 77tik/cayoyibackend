/*
 * @Author: harvey_li
 * @Date: 2022-04-12 18:05:13
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-08-05 17:13:14
 * @FilePath: \frontend\src\utils\flexible.ts
 * @Description:
 */
import store from '../store';
import Operator from './Operator';

((window, store, Operator) => {
  let screenRatioByDesign = 16 / 9; // 设计稿宽高比
  // 延时时间过长时，使用F11切换屏幕全屏到非全屏，页面尺寸来不及进行重置
  let delay = 50; // 防抖延时ms
  let minWidth = 1200; // 最小宽度
  let minHeight = minWidth / screenRatioByDesign;
  let grids = 1920; // 页面栅格份数
  let designWidth = 1920;
  let docEle = document.documentElement;
  docEle.style.minWidth = `${designWidth}rem`;
  docEle.style.minHeight = `${designWidth / screenRatioByDesign}rem`;

  const setHtmlFontSize = () => {
    const clientWidth =
      docEle.clientWidth > minWidth ? docEle.clientWidth : minWidth;
    const clientHeight =
      docEle.clientHeight > minHeight ? docEle.clientHeight : minHeight;
    let screenRatio = clientWidth / clientHeight;

    let fontSize =
      ((screenRatio > screenRatioByDesign
        ? screenRatioByDesign / screenRatio
        : 1) *
        clientWidth) /
      grids;

    docEle.style.fontSize = fontSize.toFixed(6) + 'px';

    store.updateFontSize(+fontSize.toFixed(6));
  };

  const setWidthRate = () => {
    const clientWidth =
      docEle.clientWidth > minWidth ? docEle.clientWidth : minWidth;
    const widthRate = clientWidth / designWidth;

    store.updateWidthRate(+widthRate.toFixed(6));
  };

  setHtmlFontSize();
  setWidthRate();

  let operator1 = new Operator('set_root_font_size');
  let operator2 = new Operator('set_width_rate');

  window.addEventListener('resize', () => {
    operator1.debounce(setHtmlFontSize, delay);
    operator2.debounce(setWidthRate);
  });
  window.addEventListener(
    'pageshow',
    (e) => {
      if (e.persisted) {
        // 浏览器后退的时候重新计算
        operator1.debounce(setHtmlFontSize, delay);
        operator2.debounce(setWidthRate);
      }
    },
    false
  );
})(window, store, Operator);
