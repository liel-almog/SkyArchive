import { useMutation, useQueryClient } from "@tanstack/react-query";
import { fileService } from "../../../../services/file.service";
import { filesKeys } from "../../queries";
import { File } from "../../../../models/file.model";
import { App } from "antd";

export interface UseDeleteMutationProps {
  id: string;
}

export const useDeleteMutation = ({ id }: UseDeleteMutationProps) => {
  const { message } = App.useApp();

  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: fileService.deleteFile,
    mutationKey: ["delete", id],
    onSettled: async () => {
      queryClient.invalidateQueries({
        queryKey: filesKeys.getFilesKey(),
      });
    },
    onMutate: async () => {
      const previousFiles = queryClient.getQueryData<File[]>(filesKeys.getFilesKey());

      const newFiles = previousFiles?.filter((file) => file.fileId.toString() !== id);
      queryClient.setQueryData(filesKeys.getFilesKey(), newFiles ?? []);

      return { previousFiles };
    },
    onError: (_error, _variables, context) => {
      if (context) {
        queryClient.setQueryData(filesKeys.getFilesKey(), context.previousFiles);
      }

      message.error({
        content: "Error deleting file",
        key: "delete-file",
      });
    },
  });

  return {
    mutation,
  };
};
