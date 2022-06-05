import React, { createContext, useContext, useEffect, useState } from "react";
import { FirebaseOptions, initializeApp } from "firebase/app";
import { User } from "../interfaces/user.interface";
import { AuthService } from "../services/auth/AuthService";
import { NoAuthRedirectException } from "../services/auth/NoAuthRedirectException";
import { isDevEnvironment } from "../utils/isDevEnvironment";
import { useRouter } from "next/router";

const AuthContext = createContext<User | null>(null);
const AuthServiceContext = createContext<AuthService | null>(null);

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

export function useAuthService(): AuthService | null {
  return useContext(AuthServiceContext);
}

function getFirebaseAppConfig(): FirebaseOptions {
  const projectId = process.env.FIREBASE_PROJECT_ID;

  if (isDevEnvironment()) {
    return {
      projectId,
      authDomain: "localhost",
      apiKey: "fake-api-key",
    };
  }

  return {
    projectId,
    apiKey: process.env.FIREBASE_API_KEY,
    authDomain: process.env.FIREBASE_AUTH_DOMAIN,
    appId: process.env.FIREBASE_APP_ID,
  };
}

interface AuthProviderProps {
  children: React.ReactNode;
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const config = getFirebaseAppConfig();
  const firebaseApp = initializeApp(config);
  const [authService] = useState(new AuthService(firebaseApp));
  const [auth, setAuth] = useState<User | null>(null);

  useEffect(() => {
    authService
      .handleAuthRedirect()
      .then(async () => {
        const user = await authService.getSignedInUser();
        setAuth(user);
      })
      .catch(async (err) => {
        if (err instanceof NoAuthRedirectException) {
          const user = await authService.getSignedInUser();
          return setAuth(user);
        }
        console.error(err);
      })
      .catch((err) => {
        console.error("catch2", err);
      });
  }, []);

  return (
    <AuthContext.Provider value={auth}>
      <AuthServiceContext.Provider value={authService}>
        {children}
      </AuthServiceContext.Provider>
    </AuthContext.Provider>
  );
}
