import React from 'react';
import { ITimeRender } from '../../../../../../models/types';
import './index.less';
import { translateDate } from '../../../../../../utils/translateDate';

const TimeRender = ({ startTime, endTime }: ITimeRender) => {
  return (
    <span>{translateDate(startTime, endTime)}</span>
  );
};

export default TimeRender;
