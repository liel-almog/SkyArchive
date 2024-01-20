import { z } from "zod";

export const loginSchema = z.object({
  password: z.string().min(8).max(255),
  email: z.string().email().min(8).max(255),
});

export type Login = z.infer<typeof loginSchema>;

export const registerSchema = z.object({
  username: z.string().min(8).max(255),
  password: z.string().min(8).max(255),
  email: z.string().email().min(8).max(255),
});

export type Register = z.infer<typeof registerSchema>;
