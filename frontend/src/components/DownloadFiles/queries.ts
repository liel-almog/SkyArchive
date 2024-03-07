import { useQuery } from "@tanstack/react-query";
import { fileService } from "../../services/file.service";

export const filesKeys = {
  getFilesKey() {
    return ["files"] as const;
  },
  updateFavoriteKey(fileId: number) {
    return ["updateFavorite", fileId] as const;
  },
};

export const useGetUserFiles = () => {
  return useQuery({
    queryKey: filesKeys.getFilesKey(),
    queryFn: fileService.getFiles,
  });
};
