import { useContext } from "react";
import { AuthServiceContext } from "../context/AuthProvider/AuthProvider";
import { AuthService } from "../services/auth/AuthService";

export function useAuthService(): AuthService | null {
  return useContext(AuthServiceContext);
}
