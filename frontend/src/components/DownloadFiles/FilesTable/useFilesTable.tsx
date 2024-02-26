import { getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useState } from "react";
import { File } from "../../../models/file.model";
import { columns } from "./filesTable.util";

export interface UseFilesTableProps {
  files: File[];
}

export const useFilesTable = ({ files }: UseFilesTableProps) => {
  const [data] = useState(() => [...files]);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return { table };
};
