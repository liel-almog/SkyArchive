import { PropsWithChildren, createContext } from "react";

export interface AuthContextProps {
  token: string;
  authenticated: boolean;
}

export const AuthContext = createContext<AuthContextProps>({
  token: "",
  authenticated: false,
});

export interface AuthProviderProps {}

export const AuthProvider = ({ children }: PropsWithChildren) => {
  return (
    <AuthContext.Provider value={{ token: "", authenticated: false }}>
      {children}
    </AuthContext.Provider>
  );
};
