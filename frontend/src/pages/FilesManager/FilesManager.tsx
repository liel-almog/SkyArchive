import { Tabs } from "antd";
import classes from "./files-manager.module.scss";
import { UploadFiles } from "../../components/UploadFiles";
import { DownloadFiles } from "../../components/DownloadFiles";

export interface FilesManagerProps {}

export const FilesManager = () => {
  return (
    <div className={classes.container}>
      <Tabs
        className={classes.tabs}
        centered
        items={[
          {
            key: "1",
            label: "העלאת קבצים",
            children: <UploadFiles />,
          },
          {
            key: "2",
            label: "הורדת קבצים",
            children: <DownloadFiles />,
          },
        ]}
      />
    </div>
  );
};
