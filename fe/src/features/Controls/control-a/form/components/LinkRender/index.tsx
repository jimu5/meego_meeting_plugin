import React, { useMemo } from 'react';
import { IconLink } from '@douyinfe/semi-icons';
import { Typography, Image } from '@douyinfe/semi-ui';
import { LinkType} from '../../../../../../models/types'
import icon_file_lark from '../../../../../../assets/icon_file_lark.svg'
import './index.less';

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
  // TODO: 根据type显示不同的icon
  return link ? (
    <Typography.Text className='link-container' link={{ href: link, target: '_blank' }}>
      {iconMap[type]}
      <span className='link-text'>查看</span>
    </Typography.Text>
    ): <span>-</span>;;
};

export default LinkRender;
