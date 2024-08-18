import { useEffect, useState } from 'react';
import { Context, BriefField, WorkItem } from '@lark-project/js-sdk';
import { Toast } from '@douyinfe/semi-ui';
import SDKClient from '@lark-project/js-sdk';
import { sdkManager } from '../../../utils';
import useSdkContext from '../../../hooks/useSdkContext';
import { getToken } from '../../../api';

interface Document {
  key: string;
  value: string;
}
const useModel = () => {
  const [sdk, setSdk] = useState<SDKClient>();
  const [fieldList, setFieldList] = useState<BriefField[]>([]);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [workItem, setWorkItem] = useState<WorkItem>();
  const content: Context | undefined = useSdkContext();

  useEffect(() => {
    var script = document.createElement('script');
    script.src = 'https://lf1-cdn-tos.bytegoofy.com/goofy/lark/op/h5-js-sdk-1.5.16.js';
    document.head.appendChild(script);
    // getToken().then(res => {
    //   console.log(res, 'res------token--');
    // });
    script.onload = () => {
      const login_info = '{{ login_info }}';
      console.log('login info: ', login_info);
      if (login_info === 'alreadyLogin') {
        const user_info = JSON.parse('{{ user_info | tojson | safe }}');
        console.log('user: ', user_info.name);
      } else {
        // 通过ready接口确认环境准备就绪后才能调用API
        window.h5sdk.ready(() => {
          console.log('window.h5sdk.ready');
          console.log('url:', window.location.href);
          // 调用JSAPI tt.requestAuthCode 获取 authorization code
          window.tt.requestAuthCode({
            appId: 'cli_a5fa909c3a71900b',
            // 获取成功后的回调
            success(res) {
              console.log('getAuthCode succeed: ', res);
              //authorization code 存储在 res.code
              // 此处通过fetch把code传递给接入方服务端Route: callback，并获得user_info
              // 服务端Route: callback的具体内容请参阅服务端模块server.py的callback()函数
              // fetch(`/callback?code=${res.code}`)
              //   .then(response2 =>
              //     response2.json().then(res2 => {
              //       console.log('getUserInfo succeed: ', res2);
              //       // 示例Demo中单独定义的函数showUser，用于将用户信息展示在前端页面上
              //     }),
              //   )
              //   .catch(function (e) {
              //     console.error(e);
              //   });
            },
            // 获取失败后的回调
            fail(err) {
              console.log(`getAuthCode failed, err:`, JSON.stringify(err));
            },
          });
        });
      }
      console.log(document, 'document----');
    };
  }, []);
  useEffect(() => {
    (async () => {
      const sdk = await sdkManager.getSdkInstance();
      setSdk(sdk);
    })();
  }, []);

  const fetchDocuments = async (fields: BriefField[], workItem: WorkItem) => {
    let result: Document[] = [];
    const requests = fields.map(field => workItem.getFieldValue(field.id));
    await Promise.all(requests).then(results => {
      for (let i = 0; i < fields.length; i++) {
        result.push({
          key: fields[i].name,
          value: results[i],
        });
      }
    });
    setDocuments(result);
  };

  useEffect(() => {
    if (!sdk || !content?.activeWorkItem) return;
    const { spaceId, workObjectId, id } = content?.activeWorkItem;
    (async () => {
      try {
        const detail = await sdk.WorkItem.load({
          spaceId,
          workObjectId,
          workItemId: id,
        });
        const data = await sdk.WorkObject.load({
          spaceId,
          workObjectId,
        });
        setWorkItem(detail);
        setFieldList(data.fieldList ?? []);
      } catch (e) {
        Toast.error(e.message);
      }
    })();
  }, [content?.activeWorkItem]);

  useEffect(() => {
    if (!Array.isArray(fieldList) || fieldList.length === 0 || !workItem) {
      return;
    }
    fetchDocuments(
      fieldList.filter(field => field.type === 'link'),
      workItem,
    );
  }, [fieldList, workItem]);

  return { documents };
};

export default useModel;
