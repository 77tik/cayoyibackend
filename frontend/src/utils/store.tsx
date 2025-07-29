/* Copyright (C) 2016-present, Yuansuan.cn */
import { createContext, useContext } from 'react';
export function createStore<T extends (...args: any) => any>(
  useExternalStore: T
) {
  // @ts-ignore
  const Context = createContext<ReturnType<T>>(null);

  function Provider({ children }) {
    const store = useExternalStore();
    return <Context.Provider value={store}>{children}</Context.Provider>;
  }
  return {
    Provider,
    Context,
    useStore: function useStore() {
      return useContext(Context);
    }
  };
}

// 界面下创建store.ts用来存储同一个父组件的子组件公用的一个store。并非全局store，会在路由切换，界面销毁后重新创建
/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
/**
 * step1

import { createStore } from '@/utils/store';
import { useLocalStore } from 'mobx-react-lite';

export function useModel() {
  const store = useLocalStore(
    (): {
      xxx:string
    } => ({
      xxx:''
    })
  );
  return store;
}
const store = createStore(useModel);
export const Provider = store.Provider;
export const useStore = store.useStore;


*/

// ************************
/**
 * step2
// 父组件使用step1👆🏻产生的store和provide使用

import {useStore, Provider} ./store
const Component = observer(()=>{
  const store = useStore();
  // 或者const {xxx,xxx,xxx} = useStore();
  return <>子组件</>
})
export default ()=>{
  return <Provider><Component/></Provider>
}

*/
