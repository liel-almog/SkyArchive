import { Login, Signup, tokenSchema } from "../models/auth.model";
import { axiosInstance } from "./index.service";

const PREFIX = "auth";

export class AuthService {
  async login({ email, password }: Login) {
    try {
      const res = await axiosInstance.post(`/${PREFIX}/login`, {
        email,
        password,
      });

      return tokenSchema.parse(res.data);
    } catch (error) {
      throw new Error("Error while logging in");
    }
  }

  async signup({ email, password, username }: Signup) {
    try {
      const res = await axiosInstance.post(`/${PREFIX}/signup`, {
        email,
        password,
        username,
      });

      return tokenSchema.parse(res.data);
    } catch (error) {
      throw new Error("Error while signing up");
    }
  }
}

export const authService = new AuthService();
