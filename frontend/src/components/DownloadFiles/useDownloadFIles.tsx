import { App } from "antd";
import { useEffect } from "react";
import { useAuthContext } from "../../context/AuthContext/useAuthProvider";
import { useGetUserFiles } from "./queries";

export const useGetFiles = () => {
  const { message } = App.useApp();
  const { user } = useAuthContext();
  if (!user) {
    throw new Error("User not found");
  }

  const query = useGetUserFiles();
  const { isError, error } = query;

  useEffect(() => {
    if (isError && error) {
      message.error({
        content: error.message,
        key: "get-user-files",
      });
    }
  }, [error, isError, message]);

  return { query };
};
