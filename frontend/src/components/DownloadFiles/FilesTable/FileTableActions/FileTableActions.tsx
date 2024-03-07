import { faStar as faRegularStar } from "@fortawesome/free-regular-svg-icons";
import {
  faDownload,
  faEdit,
  faStar as faSolidStar,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { CellContext } from "@tanstack/react-table";
import { Button } from "antd";
import { useState } from "react";
import { File } from "../../../../models/file.model";
import { fileService } from "../../../../services/file.service";
import { FileRenameActionModal } from "./FileRenameActionModal";
import classes from "./file-table-actions.module.scss";
import { useFavoriteMutation } from "./useFavoriteMutation";

export interface FileTableActionsProps {
  info: CellContext<File, unknown>;
}

export const FileTableActions = ({ info }: FileTableActionsProps) => {
  const [isRenameModelOpen, setIsRenameModalOpen] = useState(false);
  const { mutation: favoriteMutation } = useFavoriteMutation();

  const starIcon = info.row.original.favorite ? faSolidStar : faRegularStar;
  return (
    <>
      <section className={classes.rowActions}>
        <Button
          onClick={() => fileService.downloadFile(info.row.original.fileId)}
          icon={<FontAwesomeIcon icon={faDownload} />}
          shape="circle"
          type="text"
        />
        <Button
          onClick={() => {
            favoriteMutation.mutate({
              id: info.row.original.fileId,
              favorite: !info.row.original.favorite,
            });
          }}
          icon={<FontAwesomeIcon icon={starIcon} />}
          shape="circle"
          type="text"
        />
        <Button
          onClick={() => setIsRenameModalOpen(true)}
          icon={<FontAwesomeIcon icon={faEdit} />}
          shape="circle"
          type="text"
        />
        <Button
          onClick={() => fileService.downloadFile(info.row.original.fileId)}
          icon={<FontAwesomeIcon icon={faTrash} />}
          shape="circle"
          type="text"
        />
      </section>
      <FileRenameActionModal
        originalDisplayName={info.row.original.displayName}
        fileId={info.row.original.fileId}
        isRenameModelOpen={isRenameModelOpen}
        setIsRenameModalOpen={setIsRenameModalOpen}
      />
    </>
  );
};
