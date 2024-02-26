import { createColumnHelper } from "@tanstack/react-table";
import { File } from "../../../models/file.model";

const columnHelper = createColumnHelper<File>();

export const columns = [
  columnHelper.accessor("displayName", {
    header: "שם קובץ",
    cell: (info) => info.getValue(),
    size: 100,
  }),
  columnHelper.accessor("uploadedAt", {
    header: "תאריך העלאה",
    cell: (info) => info.getValue().toISOString(),
    size: 100,
  }),
  columnHelper.accessor("favorite", {
    header: "מועדף",
    cell: (info) => (info.getValue() ? "כן" : "לא"),
  }),
  columnHelper.accessor("status", {
    header: "סטטוס",
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor("size", {
    header: "גודל",
    cell: (info) => info.getValue(),
  }),
];
