import { PropsWithChildren } from "react";
import { AntdConfigProvider } from "./AntdConfigProvider";
import { AntdAppProvider } from "./AntdAppProvider";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { AuthProvider } from "./AuthContext/AuthProvider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export const GlobalProvider = ({ children }: PropsWithChildren) => {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <AntdConfigProvider>
          <AntdAppProvider>{children}</AntdAppProvider>
        </AntdConfigProvider>
      </AuthProvider>
      {import.meta.env.DEV && <ReactQueryDevtools />}
    </QueryClientProvider>
  );
};
