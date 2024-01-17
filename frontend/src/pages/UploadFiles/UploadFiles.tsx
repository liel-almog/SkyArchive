import { faArrowUpFromBracket, faFileUpload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Upload, UploadProps } from "antd";
import { uploadService } from "../../services/upload.service";
import classes from "./upload-files.module.scss";
const { Dragger } = Upload;

export interface UploadFilesProps {}

export const UploadFiles = () => {
  const handleCustomRequest: UploadProps["customRequest"] = async ({
    file,
    onSuccess,
    onError,
  }) => {
    try {
      let uploadFile: File;
      console.log(file instanceof File, "is file");
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

  return (
    <section className={classes.container}>
      <h2>Upload your files</h2>
      <Dragger customRequest={handleCustomRequest} className={classes.uploadingSection} multiple>
        <p className="ant-upload-drag-icon">
          <FontAwesomeIcon color="grey" size="3x" icon={faArrowUpFromBracket} />
        </p>
        <article className="ant-upload-text">
          <p>Drag 'n' drop some files here, or click to select files</p>
          <br />
        </article>
        <Button size="large" icon={<FontAwesomeIcon size="lg" icon={faFileUpload} />}>
          Upload
        </Button>
      </Dragger>
    </section>
  );
};
