import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import {
  SignupWithConfirm,
  signupWithConfirmSchema,
} from "../../models/auth.model";

export const useSignup = () => {
  const methods = useForm<SignupWithConfirm>({
    defaultValues: {
      email: "",
      password: "",
      username: "",
      confirmPassword: "",
    },
    resolver: zodResolver(signupWithConfirmSchema),
  });

  const onSubmit: SubmitHandler<SignupWithConfirm> = (data) => {
    console.log(data);
  };

  return {
    methods,
    onSubmit,
  };
};
