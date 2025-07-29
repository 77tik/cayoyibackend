/*
 * @Author: harvey_li
 * @Date: 2022-08-09 09:23:55
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-09-02 13:20:10
 * @FilePath: \frontend\src\app.ts
 * @Description:
 */
// 运行时配置

// 全局初始化数据配置，用于 Layout 用户信息和权限初始化
// 更多信息见文档：https://next.umijs.org/docs/api/runtime-config#getinitialstate

// 引入styled 全局声明
import '@/utils/declareStyledGlobal';
import '@/utils/flexible';
export async function getInitialState(): Promise<{ name: string }> {
  return { name: '@umijs/max' };
}

export const layout = () => {
  return {
    logo: 'https://img.alicdn.com/tfs/TB1YHEpwUT1gK0jSZFhXXaAtVXa-28-27.svg',
    menu: {
      locale: false,
    },
  };
};
