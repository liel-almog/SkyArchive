import { SubmitHandler, useForm } from "react-hook-form";
import { Login, loginSchema } from "../../models/auth.model";
import { zodResolver } from "@hookform/resolvers/zod";

export const useLogin = () => {
  const methods = useForm<Login>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(loginSchema),
  });

  const onSubmit: SubmitHandler<Login> = (data) => {
    console.log(data);
  };

  return {
    methods,
    onSubmit,
  };
};
