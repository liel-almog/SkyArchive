import { BlobClient } from "@azure/storage-blob";
import { startUploadSchema } from "../models/upload.model";
import { authenticatedInstance } from "./index.service";

const PREFIX = "upload" as const;

export class UploadService {
  async uploadFile(file: File) {
    const { id, signedUrl } = await this.startUpload(file);

    const blobClient = new BlobClient(signedUrl);
    const blockBlobClient = blobClient.getBlockBlobClient();
    await blockBlobClient.uploadData(file);

    await this.completeUpload(id);
  }

  private async startUpload(file: File) {
    const SIZE = file.size;
    const { data } = await authenticatedInstance.post(`/${PREFIX}/start`, {
      fileName: file.name,
      size: SIZE,
    });

    return startUploadSchema.parse(data);
  }

  private async completeUpload(id: number) {
    await authenticatedInstance.post(`/${PREFIX}/complete/${id}`);
  }
}

export const uploadService = new UploadService();
