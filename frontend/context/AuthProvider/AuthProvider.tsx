import { Alert, Snackbar } from "@mui/material";
import { FirebaseOptions, initializeApp } from "firebase/app";
import { createContext, useEffect, useState } from "react";
import { User } from "../../interfaces/user.interface";
import { AuthService } from "../../services/auth/AuthService";
import { isDevEnvironment } from "../../utils/isDevEnvironment";
import { useRedirectAuth } from "./useRedirectAuth";

export const AuthContext = createContext<User | null>(null);
export const AuthServiceContext = createContext<AuthService | null>(null);

interface AuthProviderProps {
  children: React.ReactNode;
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

export default function AuthProvider({ children }: AuthProviderProps) {
  const config = getFirebaseAppConfig();
  const firebaseApp = initializeApp(config);
  const [authService] = useState(new AuthService(firebaseApp));
  const { user, error } = useRedirectAuth(authService);
  const [errorAlertOpen, setErrorAlertOpen] = useState(false);

  useEffect(() => {
    setErrorAlertOpen(!!error);
  }, [error]);

  function handleErrorAlertClose() {
    setErrorAlertOpen(false);
  }

  return (
    <AuthContext.Provider value={user}>
      <AuthServiceContext.Provider value={authService}>
        <Snackbar
          open={errorAlertOpen}
          autoHideDuration={4000}
          onClose={handleErrorAlertClose}>
          <Alert severity="error">{error?.message}</Alert>
        </Snackbar>
        {children}
      </AuthServiceContext.Provider>
    </AuthContext.Provider>
  );
}
