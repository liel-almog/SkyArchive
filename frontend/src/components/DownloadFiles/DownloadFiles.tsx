import { FilesTable } from "./FilesTable";
import classes from "./download-files.module.scss";
import { useGetFiles } from "./useDownloadFIles";

export interface DownloadFilesProps {}

export const DownloadFiles = () => {
  const {
    query: { data, isSuccess },
  } = useGetFiles();

  return (
    <section className={classes.container}>
      <h2>הקבצים שלך</h2>
      {isSuccess ? <FilesTable files={data} /> : null}
    </section>
  );
};
