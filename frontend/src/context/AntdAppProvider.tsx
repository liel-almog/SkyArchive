import { App } from "antd";
import { PropsWithChildren } from "react";

export const AntdAppProvider = ({ children }: PropsWithChildren) => (
  <App style={{ width: "100%", height: "100%" }} message={{ duration: 2, rtl: true }}>
    {children}
  </App>
);
