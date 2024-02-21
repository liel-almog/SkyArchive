import { File } from "../../../models/file.model";
import classes from "./files-table.module.scss";

export interface FilesTableProps {
  files: File[];
}

export const FilesTable = ({ files }: FilesTableProps) => {
  return (
    <table>
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <th key={header.id}>
                {header.isPlaceholder
                  ? null
                  : flexRender(header.column.columnDef.header, header.getContext())}
              </th>
            ))}
          </tr>
        ))}
      </thead>
    </table>
  );
};
