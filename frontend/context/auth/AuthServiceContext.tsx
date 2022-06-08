import { createContext } from "react";
import { AuthService } from "../../services/auth/AuthService";

export const AuthServiceContext = createContext<AuthService | null>(null);
