import axios from 'axios';
import { PLUGIN_ID, PLUGIN_SECRET } from '../constants';

interface ResWrapper<T = {}> {
  message: string;
  statusCode: number;
  data: T;
}
axios.interceptors.request.use(
  config => {
    return config;
  },
  err => Promise.reject(err),
);
axios.interceptors.response.use(
  function (response) {
    response.data.statusCode = response.data?.status_code;
    delete response.data?.status_code;
    return response;
  },
  function (error) {
    return Promise.reject(error);
  },
);

interface AuthRes {
  code: string;
  state: string;
}
/**
 * Login authentication
 * @param data
 * @returns
 */
export function authAgree(code: string) {
  return axios.get<ResWrapper<AuthRes>>(`/login?code=${code}`).then(res => res.data);
}

export function getToken() {
  return axios
    .post<ResWrapper<AuthRes>>(`/open_api/authen/plugin_token`, {
      plugin_id: PLUGIN_ID,
      plugin_secret: PLUGIN_SECRET,
      type: 0,
    })
    .then(res => {
      const { expire_time, token } = res.data?.data ?? {};
      // 这里需要记录token的有效期，每次调用接口前如果过期需要重新请求
      return token;
    });
}

let lang = window.navigator.language;
console.log(lang);
