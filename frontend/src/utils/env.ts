import { z } from "zod";

export const envSchema = z.object({
  VITE_BACKEND_URL: z.string().url(),
});

export const validateEnv = () => {
  const parsedEnv = envSchema.safeParse(import.meta.env);
  if (!parsedEnv.success) {
    throw new Error("Invalid environment variables");
  }

  return parsedEnv.data;
};
