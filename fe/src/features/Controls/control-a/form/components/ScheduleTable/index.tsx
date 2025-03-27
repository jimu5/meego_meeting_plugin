import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Empty, Table, Toast } from '@douyinfe/semi-ui';
import { LinkType, WorkItemMeeting } from '../../../../../../models/types';
import TextRender from '../TextRender';
import TimeRender from '../TimeRender';
import LinkRender from '../LinkRender';
import StatusRender from '../StatusRender';
import ActionRender from '../ActionRender';
import './index.less';
import { useCalendarStore } from '../../calendarStore';
import { getBindMeetings } from '../../../../../../models/api';
import UserProfileComp from '../UserProfile';
import { mockData } from '../../constants';
import { sdkManager } from '../../../../../../utils';
import noContent from '../../../../../../assets/noContent.svg';

interface IScheduleTable {
  disabled: boolean;
  SDKReady: boolean;
}
const ScheduleTable = ({ disabled, SDKReady }: IScheduleTable) => {
  const store = useCalendarStore();
  const [pageNumber, setPageNumber] = useState(1);
  const [pageSize] = useState(20);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [dataSource, setDataSource] = useState<WorkItemMeeting[]>(disabled ? mockData : []); //
  const columns = useMemo(
    () => [
      {
        title: '会议名称', // 日程名称
        dataIndex: 'calendar_event_name',
        width: '150px',
        ellipsis: true,
        fixed: true,
        render: (text, record: WorkItemMeeting, index) => {
          const onclick = () => {
            sdkManager.getSdkInstance().then(sdk => {
              sdk?.navigation?.open(record.calendar_event_app_link);
            });
          };
          return (
            <div onClick={onclick}>
              <TextRender text={text} link="null" />
            </div>
          );
        },
      },
      {
        title: '组织者',
        dataIndex: 'calendar_event_organizer',
        width: '80px',
        render: (text, record: WorkItemMeeting, index) => {
          return <TextRender text={record.calendar_event_organizer?.display_name} />;
        },
      },
      {
        title: '描述',
        dataIndex: 'calendar_event_desc',
        width: '180px',
        ellipsis: true,
        render: (text, record, index) => {
          return <TextRender text={text} />;
        },
      },
      {
        title: '时间',
        dataIndex: 'meeting_time',
        width: '190px',
        render: (text, record, index) => {
          return (
            <TimeRender
              startTime={record?.meeting_time?.start_time}
              endTime={record?.meeting_time?.end_time}
            />
          );
        },
      },
      {
        title: '妙记',
        dataIndex: 'meeting_record_url',
        width: '100px',
        render: (text, record, index) => {
          return <LinkRender link={record.meeting_record_url} type={LinkType.MIAOJI} />;
        },
      },
      {
        title: '状态',
        dataIndex: 'meeting_status',
        width: '80px',
        render: (text, record, index) => {
          return <StatusRender status={record?.meeting_status} />;
        },
      },
      {
        title: '参与者',
        dataIndex: 'meeting_participant_count_accumulated',
        width: '80px',
        render: (text, record, index) => {
          return <TextRender text={text} />;
        },
      },
      {
        title: '关联日程操作人',
        dataIndex: 'bind_operator',
        width: '120px',
        render: (text, record: WorkItemMeeting, index) => {
          return (
            <UserProfileComp
              avatar={record.bind_operator_info?.avatar_url}
              name={record?.bind_operator_info?.name_cn}
              userId={record.bind_operator_info?.meego_user_key}
            />
          );
        },
      },
      {
        title: '操作',
        dataIndex: 'action',
        width: '100px',
        render: (text, record: WorkItemMeeting, index) => {
          return (
            <ActionRender
              meetingId={record?.meeting_id}
              isException={Boolean(record?.calendar_event_recurrence)}
              calendarId={record.calendar_id}
              calendarEventId={record.calendar_event_id}
            />
          );
        },
      },
    ],
    [],
  );

  // 获取日程信息
  const getData = useCallback(
    async (page: number = 1) => {
      setPageNumber(page);
      setLoading(true);
      getBindMeetings(
        {
          page_number: page,
          page_size: pageSize,
          // 空间信息
          project_key: store.projectId,
          work_item_id: store.workItemId,
          work_item_type_key: store.workItemTypeKey,
        },
        store.userId,
      )
        .then(res => {
          setDataSource(res?.meetings ?? []);
          setTotal(res?.total ?? 0);
        })
        .catch(err => {
          Toast.error('网络状况不佳，请重试！');
        })
        .finally(() => {
          setLoading(false);
          store?.switchUpdateSign(false);
        });
    },
    [pageNumber, pageSize, store.projectId, store.workItemTypeKey, store.workItemId, store.userId],
  );

  useEffect(() => {
    // 配置侧使用demo数据
    if (disabled || !SDKReady || !store.projectId) {
      return;
    }
    getData();
  }, [disabled, SDKReady, store.projectId, getData]);

  useEffect(() => {
    if (store?.updateSign && store.projectId) {
      getData();
    }
  }, [store?.updateSign, getData]);

  const handlePageChange = useCallback(
    (page: number) => {
      getData(page);
    },
    [getData],
  );

  const pagination = useMemo(
    () => ({
      currentPage: pageNumber,
      pageSize,
      total,
      onPageChange: handlePageChange,
    }),
    [pageNumber, pageSize, total, handlePageChange],
  );

  return (
    <Table
      style={{ minHeight: '300px' }}
      loading={loading}
      columns={columns}
      dataSource={dataSource}
      pagination={pagination}
      empty={<Empty title={'暂无数据'} />}
    />
  );
};

export default ScheduleTable;
