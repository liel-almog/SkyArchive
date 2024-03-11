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
import { useFavoriteMutation } from "./useFavoriteMutation";
import { useDeleteMutation } from "./useDeleteMutation";

export interface FileTableActionsProps {
  info: CellContext<File, unknown>;
}

export const FileTableActions = ({ info }: FileTableActionsProps) => {
  const [isRenameModelOpen, setIsRenameModalOpen] = useState(false);

  const { favorite, fileId, displayName } = info.row.original;
  const { mutation: favoriteMutation } = useFavoriteMutation({
    id: fileId.toString(),
  });
  const { mutation: deleteMutation } = useDeleteMutation({ id: fileId.toString() });

  const starIcon = favorite ? faSolidStar : faRegularStar;
  return (
    <>
      <section role="actions">
        <Button
          onClick={async () => {
            const { url, fileName } = await fileService.downloadFile(fileId);

            const link = document.createElement("a");

            link.download = fileName;
            link.href = url;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
          }}
          icon={<FontAwesomeIcon icon={faDownload} />}
          shape="circle"
          type="text"
        />
        <Button
          onClick={() => {
            favoriteMutation.mutate({
              id: fileId,
              favorite: !favorite,
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
          onClick={() => deleteMutation.mutate({ id: fileId })}
          icon={<FontAwesomeIcon icon={faTrash} />}
          shape="circle"
          type="text"
        />
      </section>
      <FileRenameActionModal
        originalDisplayName={displayName}
        fileId={fileId}
        isRenameModelOpen={isRenameModelOpen}
        setIsRenameModalOpen={setIsRenameModalOpen}
      />
    </>
  );
};
