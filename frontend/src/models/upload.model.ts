import { z } from "zod";

export const startUploadSchema = z.object({
  id: z.number().min(1),
});
