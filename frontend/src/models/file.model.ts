import { z } from "zod";
import { customValidation } from "../utils/zod/errorMap";

export const startUploadSchema = z.object({
  id: z.number().min(1),
  signedUrl: z.string(),
});

export const fileSchema = z.object({
  fileId: z.number().int().min(1),
  displayName: z.string(),
  uploadedAt: customValidation.dateLikeToDate,
  favorite: z.boolean(),
  status: z.string(),
  size: z.number().int().min(1),
});

export type File = z.infer<typeof fileSchema>;
