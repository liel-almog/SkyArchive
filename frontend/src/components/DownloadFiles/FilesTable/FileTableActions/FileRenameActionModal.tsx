import { Input, Modal } from "antd";
import { useState } from "react";
import { useDisplayNameMutation } from "./useDisplayNameMutation";

export interface FileRenameActionModalProps {
  originalDisplayName: string;
  fileId: number;
  isRenameModelOpen: boolean;
  setIsRenameModalOpen: (isOpen: boolean) => void;
}

export const FileRenameActionModal = ({
  originalDisplayName,
  isRenameModelOpen,
  setIsRenameModalOpen,
  fileId,
}: FileRenameActionModalProps) => {
  const [displayName, setDisplayName] = useState(originalDisplayName);
  const { mutation: displayNameMutation } = useDisplayNameMutation();

  const handleRename = () => {
    setIsRenameModalOpen(false);
    displayNameMutation.mutate({
      id: fileId,
      displayName,
    });
  };

  return (
    <Modal
      centered
      open={isRenameModelOpen}
      onOk={handleRename}
      onCancel={() => setIsRenameModalOpen(false)}
      title="Rename"
    >
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handleRename();
        }}
      >
        <Input
          value={displayName}
          onChange={(event) => {
            setDisplayName(event.target.value);
          }}
        />
      </form>
    </Modal>
  );
};
