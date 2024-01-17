import { startUploadSchema } from "../models/upload.model";
import { axiosInstance } from "./index.service";

const PREFIX = "upload" as const;

export class UploadService {
  async uploadFile(file: File) {
    const { id, CHUNKS_COUNT, CHUNK_SIZE } = await this.startUpload(file);

    const upload = (chunkIndex: number) => {
      const formData = new FormData();
      formData.append(
        "file",
        file.slice(chunkIndex * CHUNK_SIZE, chunkIndex * CHUNK_SIZE + CHUNK_SIZE)
      );
      formData.append("index", chunkIndex.toString());
      formData.append("total", CHUNKS_COUNT.toString());
      return axiosInstance.post(`/${PREFIX}/chunk/${id}/${chunkIndex}`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
    };

    for (let chunkIndex = 0; chunkIndex < CHUNKS_COUNT; chunkIndex++) {
      await upload(chunkIndex);
    }

    await this.completeUpload(id);
  }

  private async startUpload(file: File) {
    const CHUNK_SIZE = 1024 * 1024; // 1MB
    const SIZE = file.size;
    const CHUNKS_COUNT = Math.ceil(SIZE / CHUNK_SIZE);
    const { data } = await axiosInstance.post(`/${PREFIX}/chunk/start`, {
      fileName: file.name,
      size: SIZE,
    });

    const { id } = startUploadSchema.parse(data);

    return { id, CHUNK_SIZE, CHUNKS_COUNT };
  }

  private async completeUpload(id: number) {
    await axiosInstance.post(`/${PREFIX}/chunk/complete/${id}`);
  }
}

export const uploadService = new UploadService();
