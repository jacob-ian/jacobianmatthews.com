import { Alert, Snackbar } from "@mui/material";
import { FirebaseOptions, initializeApp } from "firebase/app";
import { useEffect, useState } from "react";
import { AuthService } from "../../services/auth/AuthService";
import { isDevEnvironment } from "../../utils/isDevEnvironment";
import { AuthContext } from "./AuthContext";
import { AuthServiceContext } from "./AuthServiceContext";
import { useRedirectAuth } from "./useRedirectAuth";

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
  const [errorAlertOpen, setErrorAlertOpen] = useState(false);
  const config = getFirebaseAppConfig();
  const firebaseApp = initializeApp(config);
  const [authService] = useState<AuthService>(new AuthService(firebaseApp));
  const auth = useRedirectAuth(authService);

  const { error } = auth;

  useEffect(() => {
    setErrorAlertOpen(!!error);
  }, [error]);

  function handleErrorAlertClose() {
    setErrorAlertOpen(false);
  }

  return (
    <AuthContext.Provider
      value={auth || { loading: true, user: null, error: null }}>
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
