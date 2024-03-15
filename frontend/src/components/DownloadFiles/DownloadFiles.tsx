import { Empty, Spin } from "antd";
import { FilesTable } from "./FilesTable";
import classes from "./download-files.module.scss";
import { useGetFiles } from "./useDownloadFIles";

export interface DownloadFilesProps {}

export const DownloadFiles = () => {
  const { query } = useGetFiles();
  const { isSuccess, data, isLoading } = query;

  if (isLoading) {
    return (
      <section className={classes.container}>
        <h2>הקבצים שלך</h2>
        <Spin />
      </section>
    );
  }

  if (isSuccess) {
    return (
      <section className={classes.container}>
        <h2>הקבצים שלך</h2>
        <FilesTable files={data} />
      </section>
    );
  }

  return (
    <section className={classes.container}>
      <h2>הקבצים שלך</h2>
      <Empty description="אין קבצים" />
    </section>
  );
};
