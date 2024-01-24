import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { SignupWithConfirm, signupWithConfirmSchema } from "../../models/auth.model";
import { authService } from "../../services/auth.service";
import { useMutation } from "@tanstack/react-query";
import { useAuthContext } from "../../context/AuthContext/useAuthProvider";
import { useNavigate } from "react-router-dom";

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

  const { handleLogin } = useAuthContext();
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationKey: ["signup"],
    mutationFn: authService.signup,
  });

  const onSubmit: SubmitHandler<SignupWithConfirm> = async (data) => {
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
