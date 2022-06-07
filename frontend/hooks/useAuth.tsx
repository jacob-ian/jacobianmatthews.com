import { useContext, useEffect } from "react";
import { User } from "../interfaces/user.interface";
import { useRouter } from "next/router";
import { AuthContext } from "../context/AuthProvider/AuthProvider";

interface OnFailOptions {
  redirectTo: string;
}

export function useAuth(onFailOptions?: OnFailOptions): User | null {
  const router = useRouter();
  const user = useContext(AuthContext);

  useEffect(() => {
    if (!onFailOptions) {
      return;
    }
    const { redirectTo } = onFailOptions;
    if (!user) {
      router.push(redirectTo);
    }
  }, [user]);

  return user;
}
