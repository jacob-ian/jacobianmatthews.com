import { useContext, useEffect } from "react";
import { useRouter } from "next/router";
import { AuthContext } from "../context/auth/AuthContext";
import { AuthResult } from "../context/auth/useRedirectAuth";

export interface OnFailOptions {
  redirectTo: string;
}

export function useAuth(onFailOptions?: OnFailOptions): AuthResult {
  const router = useRouter();
  const auth = useContext(AuthContext);
  const { user, loading, error } = auth;

  useEffect(() => {
    if (!onFailOptions || loading) {
      return;
    }

    const { redirectTo } = onFailOptions;
    if (!user || !!error) {
      router.push(redirectTo);
    }
  }, [user, loading]);

  return auth;
}
