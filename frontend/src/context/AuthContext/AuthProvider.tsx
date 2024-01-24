import { PropsWithChildren, createContext, useEffect, useState } from "react";
import { useLocalStorage } from "../../hooks/useLocalStorage";
import { User, userSchema } from "../../models/user.model";
import { Buffer } from "buffer";

export interface AuthContextProps {
  token: string;
  user: User | null;
  authenticated: boolean;
  handleLogin: (token: string) => void;
  handleLogout: () => void;
}

export const AuthContext = createContext<AuthContextProps>({
  token: "",
  authenticated: false,
  handleLogin: () => {},
  handleLogout: () => {},
  user: null,
});

export interface AuthProviderProps {}

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [token, setToken] = useLocalStorage("token", "");
  const [user, setUser] = useState<User | null>(null);
  const [authenticated, setAuthenticated] = useState(false);

  const parseJwt = (token: string) => {
    try {
      const user = JSON.parse(Buffer.from(token.split(".")[1], "base64").toString());
      return userSchema.parse(user);
    } catch (e) {
      return null;
    }
  };

  const handleLogin = (token: string) => {
    const user = parseJwt(token);
    if (!user) {
      return;
    }

    setUser(user);
    setToken(token);
    setAuthenticated(true);
  };

  const handleLogout = () => {
    setUser(null);
    setToken("");
    setAuthenticated(false);
  };

  useEffect(() => {
    if (token) {
      handleLogin(token);
    }
  }, [token]);

  return (
    <AuthContext.Provider
      value={{
        token,
        authenticated,
        handleLogin,
        handleLogout,
        user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
