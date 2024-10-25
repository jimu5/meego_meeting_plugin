import React, { useMemo } from 'react';
import { StatusMap } from '../../../../../../models/types';
import './index.less';
interface IStatusRender {
  status: StatusMap;
}

const StatusRender = ({ status = StatusMap.CALLED }: IStatusRender) => {
  const statusMap = useMemo(() => ({
    [StatusMap.CALLED]: {
      backgroundColor: 'rgba(143, 149, 158, 0.20)',
      label: '未开始',
    },
    [StatusMap.IN_PROGRESS]: {
      backgroundColor: 'rgba(61, 188, 47, 0.20)',
      label: '进行中',
    },
    [StatusMap.FINISH]:{
      backgroundColor: 'rgba(114, 123, 238, 0.20)',
      label: '已结束',
    }
  }), []);
  return (
    <div
      className='status-tag'
      style={{
        backgroundColor: statusMap[status].backgroundColor
      }}
    >
      {statusMap[status].label}
    </div>
  );
};

export default StatusRender;
