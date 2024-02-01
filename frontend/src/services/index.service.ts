import axios from "axios";
import { validateEnv } from "../utils/env";

const { VITE_BACKEND_URL } = validateEnv();

export const axiosInstance = axios.create({
  baseURL: `${VITE_BACKEND_URL}/api`,
});

export const authenticatedInstance = axios.create({
  baseURL: `${VITE_BACKEND_URL}/api`,
  withCredentials: true,
});

authenticatedInstance.interceptors.request.use(
  async (config) => {
    const tokenStr = localStorage.getItem("token");
    if (!tokenStr) {
      throw new Error("No token found");
    }

    const token = JSON.parse(tokenStr);

    config.headers.Authorization = `Bearer ${token}`;

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
