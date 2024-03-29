import { useMutation, useQueryClient } from "@tanstack/react-query";
import { fileService } from "../../../../services/file.service";
import { filesKeys } from "../../queries";
import { File } from "../../../../models/file.model";
import { App } from "antd";

export interface UseFavoriteMutationProps {
  id: string;
}

export const useFavoriteMutation = ({ id }: UseFavoriteMutationProps) => {
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const mutation = useMutation({
    mutationFn: fileService.updateFavorite,
    mutationKey: ["update-favorite", id],
    onMutate: async (variables) => {
      const { getFilesKey: getFiles } = filesKeys;

      await queryClient.cancelQueries({
        queryKey: getFiles(),
      });

      const previousFiles = queryClient.getQueryData<File[]>(getFiles());
      const newFiles = previousFiles?.map((file) => {
        if (file.fileId === variables.id) {
          return {
            ...file,
            favorite: variables.favorite,
          };
        }

        return file;
      });

      queryClient.setQueryData(getFiles(), newFiles ?? []);

      return { previousFiles, variables };
    },
    onError: (_error, _variables, context) => {
      if (context) {
        queryClient.setQueryData(filesKeys.getFilesKey(), context.previousFiles);
      }
      message.error({
        content: "Error updating favorite",
        key: "update-favorite",
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: filesKeys.getFilesKey(),
      });
    },
  });

  return {
    mutation,
  };
};
