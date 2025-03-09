import React from 'react';
import tableA from './table';
import formA from './form';
import CustomFieldMetaModel from './CustomFieldMetaModel';
import { CONTROL_KEY } from './constants';

import {IControlFormItemProps} from 'src/constants/type'

export const c = {
  key: CONTROL_KEY,
  renderer: {
    fieldMeta: CustomFieldMetaModel,
    render: {
      tableA,
      formA,
    },
  },
};


const App: React.FC<IControlFormItemProps> = (...props) => {
  return (
    <formA.component.display params={{props}} SDKReady={true} />
  );
}

export default App;