import React, { useState } from 'react';
import { Slider, Table } from '@douyinfe/semi-ui';
import { IControlFormItemProps } from '../../../../../constants/type';
import { columns, mockData } from '../constants';
import './index.less';

const EditFormItem = (props: IControlFormItemProps) => {
  const [value, setValue] = useState(props.value);

  return (
    <div className='container'>
      {/* <BindSchedule /> */}
      <Table columns={columns} dataSource={mockData} pagination={false} />
    </div>
  );
};

export default EditFormItem;
