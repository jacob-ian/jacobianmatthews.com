import { NextPage } from "next";
import { useRouter } from "next/router";
import AppleLoginButton from "../components/Login/SocialLoginButton/AppleLoginButton/AppleLoginButton";
import GoogleLoginButton from "../components/Login/SocialLoginButton/GoogleLoginButton/GoogleLoginButton";
import { useAuth, useAuthService } from "../hooks/useAuth";
import styles from "../styles/Login.module.scss";

const Login: NextPage = () => {
  const auth = useAuth();
  const authService = useAuthService();
  const router = useRouter();

  if (auth) {
    router.push("/dashboard");
  }

  function handleLoginButtonClick(provider: "apple" | "google") {
    if (!authService) {
      return;
    }
    if (provider === "apple") {
      authService.signInWithApple();
    }
    if (provider === "google") {
      authService.signInWithGoogle();
    }
  }

  return (
    <>
      <h1>Login</h1>
      {auth ? (
        <p>Logging in...</p>
      ) : (
        <div className={styles["login-container"]}>
          <GoogleLoginButton
            onClick={() => handleLoginButtonClick("google")}
            disabled={!authService}
          />
          <AppleLoginButton
            onClick={() => handleLoginButtonClick("apple")}
            disabled={!authService}
          />
        </div>
      )}
    </>
  );
};

export default Login;
