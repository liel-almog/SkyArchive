import { useMutation, useQueryClient } from "@tanstack/react-query";
import { fileService } from "../../../../services/file.service";
import { filesKeys } from "../../queries";
import { File } from "../../../../models/file.model";
import { App } from "antd";

export const useDisplayNameMutation = () => {
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const mutation = useMutation({
    mutationFn: fileService.updateDisplayName,
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
            displayName: variables.displayName,
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
        content: "Error updating display name",
        key: "update-display-name",
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
