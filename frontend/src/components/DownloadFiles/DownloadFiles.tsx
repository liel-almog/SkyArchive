import { FilesTable } from "./FilesTable";
import classes from "./download-files.module.scss";
import { useDownloadFiles } from "./useDownloadFIles";

export interface DownloadFilesProps {}

export const DownloadFiles = () => {
  const {
    query: { data, isSuccess },
  } = useDownloadFiles();

  return (
    <section className={classes.container}>
      <h2>הורדת קבצים</h2>
      {isSuccess ? <FilesTable files={data} /> : null}
    </section>
  );
};
