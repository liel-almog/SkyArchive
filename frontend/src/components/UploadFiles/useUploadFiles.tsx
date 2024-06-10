import { UploadProps } from "antd";
import { fileService } from "../../services/file.service";
import { useQueryClient } from "@tanstack/react-query";

export const useUploadFiles = () => {
  const queryClient = useQueryClient()
  const handleCustomRequest: UploadProps["customRequest"] = async ({
    file,
    onSuccess,
    onError,
  }) => {
    try {
      let uploadFile: File;
      if (file instanceof File) {
        uploadFile = file;
      } else {
        uploadFile = new File([new Blob([file])], "file");
      }

      await fileService.uploadFile(uploadFile);
      await queryClient.refetchQueries({
        queryKey: ["files"]
      })

      if (onSuccess) {
        onSuccess({
          status: "success",
        });
      }
    } catch (error) {
      if (error instanceof Error) {
        if (onError) {
          onError(error);
        }
      }
    }
  };

  return {
    handleCustomRequest,
  };
};
