import { Login, Register } from "../models/auth.model";

export class AuthService {
  async login({ email, password }: Login) {
    return { email, password };
  }

  async register({ email, password, username }: Register) {
    return { email, password, username };
  }
}

export const authService = new AuthService();
