import { useEffect, useState } from "react";
import { User } from "../../interfaces/user.interface";
import { AuthService } from "../../services/auth/AuthService";
import { NoAuthRedirectException } from "../../services/auth/NoAuthRedirectException";
import { UnauthenticatedException } from "../../services/http/UnauthenticatedException";

interface RedirectAuthResult {
  loading: boolean;
  user: User | null;
  error: Error | null;
}

export function useRedirectAuth(authService: AuthService): RedirectAuthResult {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    let isCancelled: boolean = false;

    async function attemptAuthentication(): Promise<User | null> {
      return authService
        .handleAuthRedirect()
        .then(() => (isCancelled ? null : authService.getSignedInUser()))
        .catch((err) => {
          if (isCancelled) {
            return null;
          }
          if (err instanceof NoAuthRedirectException) {
            return authService.getSignedInUser();
          }
          throw err;
        });
    }

    attemptAuthentication()
      .then((authUser) => {
        if (isCancelled || !authUser) {
          return;
        }
        setLoading(false);
        setUser(authUser);
      })
      .catch((err) => {
        if (isCancelled) {
          return;
        }
        setLoading(false);
        if (err instanceof UnauthenticatedException) {
          return;
        }

        // TODO: Remove
        console.error(err);
        setError(
          new Error(
            "An error occurred while signing you in. Please try again.",
          ),
        );
      });

    return () => {
      isCancelled = true;
    };
  }, []);

  return { user, error, loading };
}
