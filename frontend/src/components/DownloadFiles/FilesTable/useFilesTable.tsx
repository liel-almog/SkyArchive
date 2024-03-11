import { getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { File } from "../../../models/file.model";
import { columns } from "./filesTable.util";

export interface UseFilesTableProps {
  files: File[];
}

export const useFilesTable = ({ files }: UseFilesTableProps) => {
  const table = useReactTable<File>({
    data: files,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return { table };
};
