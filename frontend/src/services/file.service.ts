import { BlobClient } from "@azure/storage-blob";
import { fileDownloadSchema, fileSchema, startUploadSchema } from "../models/file.model";
import { authenticatedInstance } from "./index.service";

const PREFIX = "file" as const;

export class FileService {
  async uploadFile(file: File) {
    const { id, signedUrl } = await this.startFileUpload(file);

    const blobClient = new BlobClient(signedUrl);
    const blockBlobClient = blobClient.getBlockBlobClient();
    const res = await blockBlobClient.uploadData(file);

    if (res._response.status !== 201) {
      throw new Error("Error uploading file");
    }

    await this.completeFileUpload(id);

    return { id };
  }

  async downloadFile(id: number) {
    const { data } = await authenticatedInstance.get(`/${PREFIX}/download/${id}`);
    const { fileName, signedUrl } = fileDownloadSchema.parse(data);
    const blobClient = new BlobClient(signedUrl);
    const blobDownloadRes = await blobClient.download();
    const blob = await blobDownloadRes.blobBody;

    if (!blob) {
      throw new Error("Error downloading file");
    }

    const url = window.URL.createObjectURL(blob);

    return { url, fileName };
  }

  async deleteFile({ id }: { id: number }) {
    await authenticatedInstance.delete(`/${PREFIX}/${id}`);
  }

  async updateDisplayName({ id, displayName }: { id: number; displayName: string }) {
    await authenticatedInstance.patch(`/${PREFIX}/${id}/display-name`, { displayName });
  }

  async updateFavorite({ id, favorite }: { id: number; favorite: boolean }) {
    await authenticatedInstance.patch(`/${PREFIX}/${id}/favorite`, { favorite });
  }

  private async startFileUpload(file: File) {
    const SIZE = file.size;
    const { data } = await authenticatedInstance.post(`/${PREFIX}/upload/start`, {
      fileName: file.name,
      size: SIZE,
      mimeType: file.type,
    });

    return startUploadSchema.parse(data);
  }

  private async completeFileUpload(id: number) {
    await authenticatedInstance.post(`/${PREFIX}/upload/complete/${id}`);
  }

  async getFiles() {
    const { data } = await authenticatedInstance.get(`/${PREFIX}`);

    return fileSchema.array().parse(data);
  }
}

export const fileService = new FileService();
