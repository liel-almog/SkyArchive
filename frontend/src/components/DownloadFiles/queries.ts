import { useQuery } from "@tanstack/react-query";
import { fileService } from "../../services/file.service";

export const downloadKeys = {
  getFiles(userId: number) {
    return ["files", userId] as const;
  },
};

export const useGetUserFiles = (userId: number) => {
  return useQuery({
    queryKey: downloadKeys.getFiles(userId),
    queryFn: fileService.getFiles,
  });
};
