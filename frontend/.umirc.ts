/*
 * @Author: harvey_li
 * @Date: 2022-07-13 13:30:21
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-08-05 16:51:46
 * @FilePath: \frontend\.umirc.ts
 * @Description:
 */
import { defineConfig } from '@umijs/max';
import px2rem from 'postcss-pxtorem';
import proxy from './config/proxy';
import routes from './config/routes';
import theme from './config/theme';

export default defineConfig({
  antd: {},
  access: {},
  model: {},
  hash: true,
  initialState: {},
  request: {},
  layout: {
    title: '@umijs/max',
  },
  proxy,
  theme,
  extraPostCSSPlugins: [
    px2rem({
      rootValue: 1,
      unitPrecision: 6,
      minPixelValue: 3, //设置要替换的最小像素值(3px会被转rem)。 默认 0
      exclude: /node_modules/i,
      mediaQuery: true,
      propList: ['*', '!border*'],
    }),
  ],
  deadCode: {},
  devtool:
    process.env.NODE_ENV === 'development' ? 'cheap-module-source-map' : false,
  // @ts-ignore
  extraBabelPlugins:
    process.env.NODE_ENV !== 'development'
      ? [['transform-remove-console', { exclude: ['error', 'warn'] }]]
      : [],
  inlineLimit: 0,
  jsMinifierOptions: {
    minifyWhitespace: true,
    minifyIdentifiers: true,
    minifySyntax: true,
    // 添加目标环境为 ES2020
    target: ['es2020'],
  },
  cssMinifierOptions: {
    minifyWhitespace: true,
    minifySyntax: true,
  },
  svgo: {},
  copy: [
    {
      from: 'public',
      to: 'dist/public',
    },
  ],
  routes,
  npmClient: 'yarn',
});
