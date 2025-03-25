import React, { useMemo } from 'react';
import { IconLink } from '@douyinfe/semi-icons';
import { Typography, Image } from '@douyinfe/semi-ui';
import { LinkType} from '../../../../../../models/types'
import icon_file_lark from '../../../../../../assets/icon_file_lark.svg'
import './index.less';
import { sdkManager } from '../../../../../../utils';

interface ILinkRender {
  link: string;
  type: LinkType;
}

const LinkRender = ({ link, type }: ILinkRender) => {
  const iconMap = useMemo(() => ({
    [LinkType.MIAOJI]: (
      <Image
        className='link-icon'
        src={icon_file_lark}
        preview={false}
      />
    ),
    [LinkType.DOC]: <IconLink className='link-icon' />
  }), [type]);

  const onClick = () => {
    sdkManager.getSdkInstance().then((sdk) => {
      sdk?.navigation?.open(link);
    });
  }

  // TODO: 根据type显示不同的icon
  return link ? (
    <Typography.Text className='link-container' link={{}}>
      {iconMap[type]}
      <span className='link-text' onClick={onClick}>查看</span>
    </Typography.Text>
    ): <span>-</span>;;
};

export default LinkRender;
