 

import axios from 'axios';
import { message } from 'antd';
import { history, useRouteMatch } from 'umi'; 
 
const JWT_TOKEN_KEY = 'JWT_TOKEN_KEY'
export const url = window.location.origin;
const CancelToken = axios.CancelToken;
axios.defaults.timeout = 1000 * 60 * 10; //380000;
const http = axios.create({
  baseURL: `${url}/api/`,
});
message.config({
  maxCount: 1,
});

// 需要登录的状态码
//key:后端返回状态码
//value:是否使用后端返回message信息
const needLoginStatus = new Map([
  [11002, false],
  [11003, true], //ip变化，使用后端返回mes信息
  [11009, true], //登录失效
]);

http.interceptors.request.use((config): any => {
  const token =
    history.location.query!.token || localStorage.getItem(JWT_TOKEN_KEY);
  if (token) {
    config.headers['Authorization'] = `${token}`;
  } else {
    history.replace('/login');
  }
  const url = config.url; 

  return config;
});

http.interceptors.response.use(
  (res): any => {
    if (res.data.code == 0) {
      return res.data;
    }

    if (needLoginStatus.has(res.data.code)) {
      if (needLoginStatus.get(res.data.code)) {
        message.error(res.data.msg);
      }
      history.replace('/login');
      return Promise.reject(res.data.msg);
    }
    message.error(res.data.msg);
    return Promise.reject(res.data.msg);
  },
  (err) => {
 
    if (err.response && err.response.status) {
      if (err.response.status === 401) {
        history.replace('/logint');
      }
      if (err.response.status === 500) {
        message.error('服务器异常');
      }
    }
    return Promise.reject(err);
  },
);

export default http;
