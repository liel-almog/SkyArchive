import axios from "axios";
import { validateEnv } from "../utils/env";

const { VITE_BACKEND_URL } = validateEnv();

export const axiosInstance = axios.create({
  baseURL: `${VITE_BACKEND_URL}/api`,
});
