import { faArrowUpFromBracket, faFileUpload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Upload, UploadProps } from "antd";
import { uploadService } from "../../services/upload.service";
import classes from "./upload-files.module.scss";
import clsx from "clsx";
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
      <h2>העלאת קבצים</h2>
      <Dragger
        customRequest={handleCustomRequest}
        className={clsx(classes.uploadingSection)}
        maxCount={7}
        multiple
        listType="picture"
      >
        <p className="ant-upload-drag-icon">
          <FontAwesomeIcon color="grey" size="2x" icon={faArrowUpFromBracket} />
        </p>
        <article className="ant-upload-text">
          <p>גרור ושחרר קבצים כאן, או לחץ כדי לבחור קבצים</p>
          <br />
        </article>
        <Button
          className={classes.uploadBtn}
          size="large"
          icon={<FontAwesomeIcon size="lg" icon={faFileUpload} />}
        >
          העלאה
        </Button>
      </Dragger>
    </section>
  );
};
