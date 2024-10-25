import React, { useState, useEffect } from 'react';
import SDK from "@lark-project/js-sdk";
import { PLUGIN_ID } from '../../constants';

let ready = false;
const readyCallbacks: Array<(next: boolean) => void> = [];

const sdk = new SDK();

sdk.config({
  pluginId: PLUGIN_ID,
  isDebug: false,
}).then(() => {
  ready = true;
  readyCallbacks.forEach(cb => cb(ready));
});

export const withJSSDKReady = (Component: React.ComponentType) => {
  return (props: any) => {
    const [SDKReady, setSDKReady] = useState(ready);
    useEffect(() => {
      if (ready) {
        setSDKReady(ready);
        return;
      }
      readyCallbacks.push((nextSDKReady: boolean) => {
        setSDKReady(nextSDKReady);
      });
    }, []);
    return <Component SDKReady={SDKReady} {...props} />
  }
};

export default sdk;
