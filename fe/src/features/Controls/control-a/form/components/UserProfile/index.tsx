import React, { lazy, useEffect, useState } from 'react';
import { Highlight, Select, Button, Tag } from '@douyinfe/semi-ui';
import DefaultAvatr from '../../../../../../assets/BG.png';
import './index.less';

/**
 * TODO: 待确认可行性
 */
// const UserProfile = lazy(
//   () =>
//     import(
//       /* webpackChunkName: "user-profile" */ '@universe-design/biz-react-lark-web/es/components/UserProfileV2'
//     ),
// );

// <UserProfile
//   // className={className}
//   // isShowCustomStatus={isShowCustomStatus}
//   // 强制刷新larkuserid，解决甘特图中已加载过的人员头像切换，人员组件不刷新的问题
//   key={'ou_e051986ab19f80d16b7b8d74f3f1235'}
//   // addContactParams={addContactParams}
//   userId={'ou_e051986ab19f80d16b7b8d74f3f1235'} //   要打开个人卡片的目标用户的 userId
//   onActionClose={() =>{}} // 关闭卡片的操作
// />
interface IUserProfileComp {
  userId: string;
  name?: string;
  avatar?: string;
}
const UserProfileComp = ({ userId, name, avatar }: IUserProfileComp) => {
  const defaultUserInfo = {
    id: 'lijie.xxx',
    avatr: DefaultAvatr,
    name: '李雷'
  };
  useEffect(() => {
    // TODO: 根据userId获取完整的用户信息
  }, []);
  return (
    <Tag avatarSrc={avatar || defaultUserInfo.avatr} avatarShape="circle" size="large">
      {name || defaultUserInfo.name}
    </Tag>
  );
};

export default UserProfileComp;
