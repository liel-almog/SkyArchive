import { SubmitHandler, useForm } from "react-hook-form";
import { Login, loginSchema } from "../../models/auth.model";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { authService } from "../../services/auth.service";

export const useLogin = () => {
  const methods = useForm<Login>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(loginSchema),
  });

  const mutation = useMutation({
    mutationKey: ["login"],
    mutationFn: authService.login,
  });

  const onSubmit: SubmitHandler<Login> = (data) => {
    mutation.mutate(data);
  };

  return {
    methods,
    onSubmit,
  };
};
