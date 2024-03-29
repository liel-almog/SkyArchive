import { z } from "zod";
import { customValidation } from "../utils/zod/errorMap";

const passwordSchema = z
  .string()
  .trim()
  .min(8)
  .regex(/^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*\W)(?!.* ).*$/, {
    message: "הסיסמה חייבת להכיל אותיות גדולות וקטנות, מספרים ותו מיוחד",
  });

export const loginSchema = z.object({
  password: passwordSchema,
  email: customValidation.email,
});

export type Login = z.infer<typeof loginSchema>;

export const signupSchema = z.object({
  username: customValidation.text.min(1).max(255),
  password: passwordSchema,
  email: customValidation.email,
});

export type Signup = z.infer<typeof signupSchema>;

export const signupWithConfirmSchema = signupSchema
  .extend({
    confirmPassword: passwordSchema,
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "הסיסמאות אינן תואמות",
    path: ["confirmPassword"],
  });

export type SignupWithConfirm = z.infer<typeof signupWithConfirmSchema>;

export const tokenSchema = z.object({
  token: z.string(),
})

export type Token = z.infer<typeof tokenSchema>;