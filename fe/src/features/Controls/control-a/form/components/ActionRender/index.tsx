import React from 'react';
import './index.less';
import { Modal, Toast } from '@douyinfe/semi-ui';
import { unbindCalendar } from '../../../../../../models/api';
import { useCalendarStore } from '../../calendarStore';

interface IActionRender {
  calendarId: string;
  calendarEventId: string;
  meetingId?: string;
  isException?: boolean;
}

const ActionRender = ({ meetingId, isException, calendarId, calendarEventId }: IActionRender) => {
  const store = useCalendarStore();
  // TODO: 根据type显示不同的icon
  const onClick = () => {
    const cb = (sign: boolean) => {
        // 1. 调用解除绑定接口
        // 2. 刷新会议列表
        unbindCalendar({
          meeting_id: meetingId,
          calendar_event_id: calendarEventId,
          calendar_id: calendarId,
          work_item_id: store.workItemId,
          work_item_type_key: store.workItemTypeKey,
          project_key: store.projectId,
          with_after_recurring_event: sign
        }, store?.userId).then((res) => {
          store.switchUpdateSign(true);
        }).catch((err) => {
          Toast.error('网络状况不佳，请重试！')
          store?.switchUpdateSign(true);
        })
    }
    if(isException) {
      Modal.warning({
        title: '是否取消关联后续重复日程？',
        content: '检测到你取消关联的是重复日程，是否取消关联后续日程？',
        okText: '是',
        cancelText: '否',
        centered: true,
        onOk: () => {
          cb(true);
        },
        onCancel: () => {
          cb(false);
        }
      });
      return;
    }
    cb(false);
  }
  return (
    <div className='action-container' onClick={onClick}>
      取消关联
    </div>
  );
};

export default ActionRender;
