import {
  SortingState,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { useState } from "react";
import { File } from "../../../models/file.model";
import { columns } from "./filesTable.util";

export interface UseFilesTableProps {
  files: File[];
}

export const useFilesTable = ({ files }: UseFilesTableProps) => {
  const [sorting, setSorting] = useState<SortingState>([]);

  const table = useReactTable<File>({
    data: files,
    columns,
    state: {
      sorting,
    },
    getCoreRowModel: getCoreRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
  });

  return { table };
};
