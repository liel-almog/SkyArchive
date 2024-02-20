import { faArrowUpFromBracket, faFileUpload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Upload } from "antd";
import clsx from "clsx";
import classes from "./upload-files.module.scss";
import { useUploadFiles } from "./useUploadFiles";
const { Dragger } = Upload;

export interface UploadFilesProps {}

export const UploadFiles = () => {
  const { handleCustomRequest } = useUploadFiles();

  return (
    <section className={classes.container}>
      <h2>העלאת קבצים</h2>
      <Dragger
        customRequest={handleCustomRequest}
        className={clsx(classes.uploadingSection)}
        // maxCount={7}
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
