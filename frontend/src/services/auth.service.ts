import { Login, Signup } from "../models/auth.model";
import { axiosInstance } from "./index.service";

const PREFIX = "auth";

export class AuthService {
  async login({ email, password }: Login) {
    return axiosInstance.post(`/${PREFIX}/login`, { email, password });
  }

  async signup({ email, password, username }: Signup) {
    return axiosInstance.post(`/${PREFIX}/signup`, {
      email,
      password,
      username,
    });
  }
}

export const authService = new AuthService();
