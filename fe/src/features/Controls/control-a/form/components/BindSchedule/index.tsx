import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Select, Button, Toast, Switch, Tooltip, Modal } from '@douyinfe/semi-ui';
import { IconInfoCircle } from '@douyinfe/semi-icons';
import classNames from 'classnames';
import { debounce } from 'lodash-es';
import { autoBindMeetings, bindCalendar, getAutoBindMeetingsStatus, searchCalendar, updateCalendar } from '../../../../../../models/api';
import { MAX_RES_LENGTH } from '../../constants';
import { CalendarItem } from '../../../../../../models/types';
import './index.less';
import { translateDateWithWeek } from '../../../../../../utils/translateDate';
import { useCalendarStore } from '../../calendarStore';

interface IBindSchedule {
  disabled: boolean;
  SDKReady: boolean;
}
const BindSchedule = ({ disabled, SDKReady }: IBindSchedule) => {
  const store = useCalendarStore();
  const [init, setInit] = useState(false);
  const [value, setValue] = useState<CalendarItem>();
  const renderOptionItem = renderProps => {
      const {
          disabled,
          selected,
          label,
          value,
          focused,
          className,
          style,
          onMouseEnter,
          onClick,
          start_time,
          end_time,
          empty,
          emptyContent,
      } = renderProps;
      const optionCls = classNames({
          ['custom-option-render']: true,
          ['custom-option-render-focused']: focused,
          ['custom-option-render-disabled']: disabled,
          ['custom-option-render-selected']: selected,
          ['custom-option-item']: true,
          ['not-started']: start_time?.timestamp ? Boolean(+new Date > start_time?.timestamp * 1000) : true
      });
      return (
          <div style={style} className={optionCls} onClick={() => onClick()} onMouseEnter={e => onMouseEnter()}>
              <div className="option-name">
                {label || "（无主题）"}
              </div>
              <div className='option-time'>{translateDateWithWeek(start_time, end_time)}</div>
          </div>
      );
  };

  const [loading, setLoading] = useState(false);
  const [bindLoading, setBindLoading] = useState(false)
  const [enable, setEnbale] = useState(false);
  const optionList: CalendarItem[] = useMemo(() => [], []);
  const [list, setList] = useState<CalendarItem[]>(optionList);

  const bindSchedleHandler = useCallback(() => {
    const _handler = () => {
      bindCalendar({
        calendar_event_id: value?.event_id ?? '',
        calendar_id: value?.organizer_calendar_id ?? '',
        work_item_id: store?.workItemId,
        work_item_type_key: store?.workItemTypeKey,
        project_key: store?.projectId
      }, store?.userId).then((res) => {
        Toast.success('关联成功');
        store?.switchUpdateSign(true);
        setValue(undefined);
      }).catch((err) => {
        Toast.error('网络状况不佳，请重试！')
      }).finally(() => {
        setValue(undefined);
      });
    }
    // 重复日程需要二次确认
    if(value?.recurrence) {
      Modal.warning({
        title: '确定关联重复日程？',
        content: '此日程为重复日程，请确认是否关联？',
        okText: '确认',
        cancelText: '取消',
        centered: true,
        onOk: () => {
          _handler();
        },
        onCancel: () => {
          setValue(undefined);
        }
      });
      return;
    }
    _handler();
  }, [store?.workItemId, store?.workItemTypeKey, store?.projectId, store?.userId, value]);

  // 默认搜索
  const handleSearch = useCallback((inputValue?: string) => {
    if(!store.userId) {
      return;
    }
    setLoading(true);
    searchCalendar({
      page_size: MAX_RES_LENGTH,
      query_word: inputValue
    }, store.userId).then(res => {
      setList((res?.items ?? [])?.map(item => ({
        ...item,
        label: item.summary,
        value: item.event_id,
      })))
      return []
    }).catch(() => {
      Toast.error('网络状况不佳，请重试！')
    }).finally(() => {
      setLoading(false);
    });
  }, [store.userId]);

  //  获取群自动关联的日程状态
  const queryAutoBindStatus = useCallback((inputValue?: string) => {
    if(!store.userId) {
      return;
    }
    setBindLoading(true)
    getAutoBindMeetingsStatus({
      project_key: store?.projectId,
      work_item_type_key: store?.workItemTypeKey,
      work_item_id: store?.workItemId,
    }, store.userId).then(res => {
      setEnbale(res.enable);
    }).finally(() => {
      setBindLoading(false)
    });
  }, [store.userId, store?.workItemId,]);


  const handleChange = newValue => {
    setValue(newValue);
  };

  // 展开下拉框请求数据
  const onDropdownVisibleChange = useCallback((vis: boolean) => {
    if(init === false) {
      handleSearch();
      setInit(true);
    }
  }, [init, handleSearch]);

  useEffect(() => {
    if(disabled) {
      return;
    }
    queryAutoBindStatus();
  }, [queryAutoBindStatus]);

  const updateCalendarHandler = () => {
    if(!store.userId) {
      return;
    }
    updateCalendar({
      work_item_id: store?.workItemId,
      work_item_type_key: store?.workItemTypeKey,
      project_key: store?.projectId
    }, store.userId).then(res => {
      store?.switchUpdateSign(true);
    }).catch(() => {
      Toast.error('网络状况不佳，请重试！')
    });
  }

  const autoBindHandler = (_enable) => {
    if(!store.userId) {
      return;
    }
    setBindLoading(true)
    autoBindMeetings({
      enable: _enable,
      work_item_id: store?.workItemId,
      work_item_type_key: store?.workItemTypeKey,
      project_key: store?.projectId
    }, store.userId).then(res => {
      setEnbale(_enable);
    }).catch(err => {
      console.error('auto bind err', err);
    }).finally( () => {
      setBindLoading(false)
    });
  }
  return (
    <div className="selection" >
      <div className='label' onClick={updateCalendarHandler}>关联日程</div>
      <Select
        className='custom-select'
        disabled={disabled}
        filter
        placeholder="请输入日程名称"
        remote
        // multiple
        showClear
        loading={loading}
        value={value}
        onChangeWithObject
        onSearch={debounce(handleSearch, 1000)}
        dropdownClassName="components-select-demo-renderOptionItem"
        optionList={list}
        onChange={handleChange}
        renderOptionItem={renderOptionItem}
        onDropdownVisibleChange={onDropdownVisibleChange}
      />
      <Button disabled={!value || disabled} theme='solid' type='primary' onClick={bindSchedleHandler}>关联</Button>
      <div className='right-switch'>
        <Tooltip className='right-tooltip' content="开启后，将当前工作项实例群添加为日程参与者并将日程分享到群内时，会自动将日程与当前实例关联。">
          <span className='right-tooltip-span'>自动关联日程</span>
          <IconInfoCircle />
        </Tooltip>
        <Switch
          className='switch'
          disabled={disabled}
          checked={enable}
          loading={bindLoading}
          size="small"
          onChange={autoBindHandler}
          aria-label="自动关联日程"
        />
      </div>
    </div>
  );
};

export default BindSchedule;
