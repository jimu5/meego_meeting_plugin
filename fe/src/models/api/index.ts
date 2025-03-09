import axios from "axios";
import { encode } from "js-base64";
// import { PLUGIN_ID, PLUGIN_SECRET } from "../../constants";
import {
  AuthRes,
  AutoBindMeetingsParams,
  BindCalendarParams,
  BindMeetingsParams,
  BindMeetingsParamsRes,
  // GetAuthCodeRes,
  GetAutoBindMeetingsStatus,
  GetAutoBindMeetingsStatusRes,
  // GetTokenAuthRes,
  ResWrapper,
  SearchCalendarParams,
  SearchCalendarRes,
  UpdateCalendarParams,
  UserPluginTokenRes,
} from "../types";
import { sdkManager } from "../../utils";
import { Toast } from "@douyinfe/semi-ui";

// 飞书项目的服务
const API_PREFIX = `https://${window.location.host}`;

// 插件自己的服务
const CUSTOM_API_PREFIX = `https://yourhost.com`;

// 飞书的授权服务
const FEISHU_AUTH_URL = "https://open.feishu.cn/open-apis/authen/v1/index";

// 飞书的开发者后台配置的重定向URL
const REDRECT_URL =
  "https://yourhost.com/api/v1/meego/lark/auth";

// 飞书的APP_ID
const FEISHU_APP_ID = "your_feishu_app_id";

export const getFeishuAuthHandler = () => {
  sdkManager.getSdkInstance().then((res) => {
    const state = document.URL;
    res?.storage?.getItem(`user_id`).then((_res) => {
      const meegoUserId = _res;
      const encodedState = encode(`${state}&meego_user_key=${meegoUserId}`);
      const openUrl = `${FEISHU_AUTH_URL}?redirect_uri=${encodeURI(
        REDRECT_URL
      )}&app_id=${FEISHU_APP_ID}&state=${encodedState}`;
      res?.navigation?.open(openUrl);
    });
  });
};

axios.interceptors.request.use(
  (config) => {
    return config;
  },
  (err) => Promise.reject(err)
);
axios.interceptors.response.use(
  function (response) {
    response.data.statusCode = response.data?.status_code;
    delete response.data?.status_code;
    return response;
  },
  function (error) {
    if (error?.response?.data?.msg) {
      Toast.error({
        content: error?.response?.data?.msg,
      });
    }

    if (error?.response?.status === 401) {
      getFeishuAuthHandler();
      return;
    }
    return Promise.reject(error);
  }
);

/**
 * Login authentication
 * @param data
 * @returns
 */
export function authAgree(code: string) {
  return axios
    .get<ResWrapper<AuthRes>>(`${API_PREFIX}/login?code=${code}`)
    .then((res) => res.data);
}

// /**
//  * 获取authcode
//  * @returns
//  */
// export function getAuthCode() {
//   return axios
//     .post<ResWrapper<GetAuthCodeRes>>(
//       `${API_PREFIX}/open_api/authen/auth_code`,
//       {
//         plugin_id: PLUGIN_ID,
//         state: "111",
//       }
//     )
//     .then((res) => {
//       const authCode = res?.data?.data?.code;
//       return authCode;
//     });
// }

// /**
//  * 获取token
//  * @returns
//  */
// export function getToken() {
//   return axios
//     .post<ResWrapper<GetTokenAuthRes>>(
//       `${API_PREFIX}/open_api/authen/plugin_token`,
//       {
//         plugin_id: PLUGIN_ID,
//         plugin_secret: PLUGIN_SECRET,
//         type: 0,
//       }
//     )
//     .then((res) => {
//       const { expire_time, token } = res.data?.data ?? {};
//       console.log(res.data, "res.data?------");
//       // TODO: 这里需要记录token的有效期，每次调用接口前如果过期需要重新请求
//       // window.localStorage.setItem(`${PLUGIN_ID}_token`, token);
//       // window.localStorage.setItem(`${PLUGIN_ID}_expire_time`, expire_time + "");
//       return res.data?.data?.token;
//     });
// }

/**
 * 用plugin token获取user plugin token
 * @returns
 */
export function getUserPluginToken(token: string, code: string) {
  return axios
    .post<ResWrapper<UserPluginTokenRes>>(
      `${API_PREFIX}/open_api/authen/user_plugin_token`,
      {
        code,
        grant_type: "authorization_code",
      },
      {
        headers: {
          "X-Plugin-Token": token,
        },
      }
    )
    .then((res) => {
      const {
        token,
        expire_time,
        refresh_token,
        refresh_token_expire_time,
        user_key,
        saas_tenant_key,
      } = res.data?.data ?? {};
      console.log("鉴权成功后的信息返回 ", res.data?.data);
      return res.data?.data ?? {};
    });
}

/**
 * 刷新token
 * @param refresh_token
 * @param token
 * @returns
 */
export function refreshToken(refresh_token: string, token: string) {
  return axios
    .post<ResWrapper<UserPluginTokenRes>>(
      `${API_PREFIX}/open_api/authen/refresh_token`,
      {
        refresh_token,
        type: 1,
      },
      {
        headers: {
          "X-Plugin-Token": token,
        },
      }
    )
    .then((res) => {
      const {
        token,
        expire_time,
        refresh_token,
        refresh_token_expire_time,
        user_key,
        saas_tenant_key,
      } = res.data?.data ?? {};
      console.log("鉴权成功后的信息返回 ", res.data?.data);
      return "";
    });
}
// let lang = window.navigator.language;
// console.log(lang);

/**
 * 搜索日程
 * @param params 搜索关键词和分页信息
 * @returns
 */
export const searchCalendar = (
  params: SearchCalendarParams,
  userId: string
) => {
  return axios
    .get<SearchCalendarRes>(
      `${CUSTOM_API_PREFIX}/api/v1/lark/calendar/search`,
      {
        params,
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
};

/**
 * 触发日程同步
 * @param params 搜索关键词和分页信息
 * @returns
 */
export const updateCalendar = (
  params: UpdateCalendarParams,
  userId: string
) => {
  return axios
    .post<Record<string, any>>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/work_item_meetings/refresh`,
      params,
      {
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
};

/**
 * 将日程绑定Meego实例
 * @param params
 * @returns
 */
export function bindCalendar(params: BindCalendarParams, userId: string) {
  return axios
    .post<SearchCalendarRes>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/calendar_event/bind`,
      params,
      {
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
}

/**
 * 将日程与Meego实例解绑
 * @param params
 * @returns
 */
export function unbindCalendar(params: BindCalendarParams, userId: string) {
  return axios
    .post<SearchCalendarRes>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/calendar_event/unbind`,
      params,
      {
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
}

/**
 * 获取绑定到meego实例上的会议详情
 * @param params
 * @returns
 */
export function getBindMeetings(params: BindMeetingsParams, userId: string) {
  return axios
    .post<BindMeetingsParamsRes>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/work_item_meetings`,
      params,
      {
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
}

/**
 * 获取绑定到meego实例上的会议详情
 * @param params
 * @returns
 */
export function autoBindMeetings(
  params: AutoBindMeetingsParams,
  userId: string
) {
  return axios
    .post<BindMeetingsParamsRes>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/work_item_meetings/chat_auto_bind`,
      params,
      {
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
}

/**
 * 获取群自动关联的日程状态
 * @param params
 * @param userId
 * @returns
 */
export function getAutoBindMeetingsStatus(
  params: GetAutoBindMeetingsStatus,
  userId: string
) {
  return axios
    .get<GetAutoBindMeetingsStatusRes>(
      `${CUSTOM_API_PREFIX}/api/v1/meego/work_item_meetings/chat_auto_bind`,
      {
        params,
        headers: {
          meego_user_key: userId,
        },
      }
    )
    .then((res) => res.data);
}
