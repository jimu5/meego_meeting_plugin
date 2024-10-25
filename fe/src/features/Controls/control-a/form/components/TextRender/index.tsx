import React from 'react';
import './index.less';
import { Typography } from '@douyinfe/semi-ui';

interface ITextRender {
  text: string;
  link?: string;
}
// TODO: 超出隐藏
const TextRender = ({ text, link }: ITextRender) => {
  return (
    <Typography.Text
        ellipsis={{
            showTooltip: true
        }}
        link={link ? { href: link, target: '_blank' } : undefined}
        style={{ width: '100%', cursor: link ? 'pointer' : 'default' }}
      >
        {text ?? '-'}
      </Typography.Text>
  );
};

export default TextRender;
