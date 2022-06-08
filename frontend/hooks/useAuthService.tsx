import { useContext } from "react";
import { AuthServiceContext } from "../context/auth/AuthServiceContext";
import { AuthService } from "../services/auth/AuthService";

export function useAuthService(): AuthService | null {
  return useContext(AuthServiceContext);
}
