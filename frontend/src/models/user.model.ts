import { z } from "zod";

export const userSchema = z.object({
  id: z.number(),
  username: z.string(),
  email: z.string(),
});

export type User = z.infer<typeof userSchema>;
