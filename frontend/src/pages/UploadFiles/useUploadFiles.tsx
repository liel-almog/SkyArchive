import { UploadProps } from "antd";
import { uploadService } from "../../services/upload.service";

export const useUploadFiles = () => {
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

      await uploadService.uploadFile(uploadFile);

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
