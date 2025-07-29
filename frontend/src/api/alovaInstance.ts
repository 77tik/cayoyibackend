import { axiosRequestAdapter } from '@alova/adapter-axios';
import { createAlova } from 'alova';
import ReactHook from 'alova/react';

// 从 proxy 配置中获取是否使用 mock
const isMock =
  process.env.NODE_ENV === 'development' &&
  require('../../config/proxy').default.dev['/api/'].useMock;

const alovaInstance = createAlova({
  // 开发环境：mock使用相对路径，非mock使用完整服务器地址
  // 生产环境：使用空字符串，让nginx处理/api路径的代理
  baseURL:
    process.env.NODE_ENV === 'development'
      ? isMock
        ? ''
        : 'http://10.0.4.66:8849'
      : '',
  statesHook: ReactHook,
  requestAdapter: axiosRequestAdapter(),
  timeout: 10000,
  // 可选的请求头等配置
  // headers: {
  //   'Authorization': 'Bearer xxx',
  // },
});

export default alovaInstance;
