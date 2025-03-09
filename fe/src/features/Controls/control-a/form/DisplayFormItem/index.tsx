import React, { useEffect } from 'react';
import { IControlFormItemProps } from '../../../../../constants/type';
import BindSchedule from '../components/BindSchedule';
import './index.less';
import ScheduleTable from '../components/ScheduleTable';
import { useCalendarStore } from '../calendarStore';

import { sdkManager } from '../../../../../utils';

const DisplayFormItem = (props: IControlFormItemProps & {SDKReady: boolean}) => {
  const [appContext, setContext] = React.useState({});
  // const [field] = useSafeFormikFieldState('field_id');
  // const disabled = props?.mode === 'configure';
  const disabled = false;
  const getContext = async () => {
    try {
      const context = await window.JSSDK.Context.load();
      setContext(context);
    } catch (e) {}
  };
  const store = useCalendarStore();
  // useEffect(() => {
  //   const authInit = async () => {
  //     const code = await getAuthCode();
  //     const token = await getToken();
  //     getUserPluginToken(token, code).then((res) => {
  //       // TODO: 使用user plugin token 做一系列的鉴权操作
  //       console.log(res, 'res------')
  //     })
  //   }
  //   authInit();
  // }, []);
  useEffect(() => {
    getContext()
    const workItemId = appContext?.activeWorkItem?.id ?? 0;
    const workItemTypeKey = appContext?.activeWorkItem?.workObjectId ?? '';
    const projectId = appContext?.mainSpace?.id ?? '';
    const userId = appContext?.loginUser?.id ?? '';
    if(userId) {
      sdkManager.getSdkInstance().then((res) => {
        res?.storage?.setItem(`user_id`, userId);
      });
    }
    if(appContext?.activeWorkItem?.id) {
      store.init({
        workItemId,
        workItemTypeKey,
        projectId,
        userId,
      });
    }
    return () => {
      store.reset();
    }
  }, [appContext?.activeWorkItem?.id, appContext?.activeWorkItem?.workObjectId, appContext?.mainSpace?.id, appContext?.loginUser?.id]);

  return (
    <div className="from_progress" >
      <BindSchedule SDKReady={props?.SDKReady} disabled={disabled} />
      <ScheduleTable SDKReady={props?.SDKReady} disabled={disabled} />
    </div>
  );
};

// export default withJSSDKReady(DisplayFormItem);
export default DisplayFormItem;
