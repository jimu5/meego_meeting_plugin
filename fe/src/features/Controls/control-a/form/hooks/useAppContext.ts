import { Context } from "@lark-project/js-sdk";
import { useEffect, useState } from "react";
import sdk from "../../../SDKHoc";

export const useAppContext = (SDKReady: boolean = false) => {
  const [context, setContext] = useState<Context | undefined>();

  useEffect(() => {
    let unwatch;
    sdk?.Context?.load?.().then((ctx) => {
      setContext(ctx);
      unwatch = ctx.watch((nextCtx) => {
        setContext(nextCtx);
      });
    });
    return () => {
      unwatch?.();
    };
  }, [SDKReady]);

  return context;
};
