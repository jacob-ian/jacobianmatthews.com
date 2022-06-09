import { createContext } from "react";
import { AuthResult } from "./useRedirectAuth";

export const AuthContext = createContext<AuthResult>({
  user: null,
  error: null,
  loading: true,
});
