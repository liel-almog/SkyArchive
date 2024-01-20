import { PropsWithChildren } from "react";
import { AntdConfigProvider } from "./AntdConfigProvider";
import { AntdAppProvider } from "./AntdAppProvider";
import { AuthProvider } from "./AuthContext/AuthProvider";

export const GlobalProvider = ({ children }: PropsWithChildren) => {
  return (
    <AuthProvider>
      <AntdConfigProvider>
        <AntdAppProvider>{children}</AntdAppProvider>
      </AntdConfigProvider>
    </AuthProvider>
  );
};
