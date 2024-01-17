import { PropsWithChildren } from "react";
import { AntdConfigProvider } from "./AntdConfigProvider";
import { AntdAppProvider } from "./AntdAppProvider";

export const GlobalProvider = ({ children }: PropsWithChildren) => {
  return (
    <AntdConfigProvider>
      <AntdAppProvider>{children}</AntdAppProvider>
    </AntdConfigProvider>
  );
};
