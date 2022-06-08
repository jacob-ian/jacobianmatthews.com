import { createContext } from "react";
import { User } from "../../interfaces/user.interface";

export const AuthContext = createContext<User | null>(null);
