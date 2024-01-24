import { SubmitHandler, useForm } from "react-hook-form";
import { Login, loginSchema } from "../../models/auth.model";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { authService } from "../../services/auth.service";
import { useAuthContext } from "../../context/AuthContext/useAuthProvider";
import { useNavigate } from "react-router-dom";

export const useLogin = () => {
  const methods = useForm<Login>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(loginSchema),
  });
  const { handleLogin } = useAuthContext();
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationKey: ["login"],
    mutationFn: authService.login,
  });

  const onSubmit: SubmitHandler<Login> = (data) => {
    mutation.mutate(data, {
      onSuccess(data) {
        handleLogin(data.token);
        navigate("/");
      },
    });
  };

  return {
    methods,
    onSubmit,
  };
};
