import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { SignupWithConfirm, signupWithConfirmSchema } from "../../models/auth.model";
import { authService } from "../../services/auth.service";
import { useMutation } from "@tanstack/react-query";

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

  const mutation = useMutation({
    mutationKey: ["signup"],
    mutationFn: authService.signup,
  });

  const onSubmit: SubmitHandler<SignupWithConfirm> = (data) => {
    mutation.mutate(data);
  };

  return {
    methods,
    onSubmit,
  };
};
