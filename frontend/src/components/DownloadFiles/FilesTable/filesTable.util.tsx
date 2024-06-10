import { createColumnHelper } from "@tanstack/react-table";
import { File } from "../../../models/file.model";
import { FileTableActions } from "./FileTableActions";
import classes from "./files-table.module.scss";

const columnHelper = createColumnHelper<File>();

function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

export const columns = [
  columnHelper.accessor("displayName", {
    header: "שם קובץ",
    cell: (info) => info.getValue(),
    size: 100,
    sortingFn: "alphanumeric",
    sortDescFirst: true,
  }),
  columnHelper.accessor("uploadedAt", {
    header: "תאריך העלאה",
    cell: (info) =>
      info
        .getValue()
        .toLocaleDateString("he-IL", { year: "numeric", month: "2-digit", day: "2-digit" }),
    size: 100,
  }),
  columnHelper.accessor("size", {
    header: "גודל",
    cell: (info) => <span className={classes.fileSize}>{formatBytes(info.getValue())}</span>,
  }),
  columnHelper.display({
    id: "actions",
    cell: (info) => <FileTableActions info={info} />,
  }),
];
