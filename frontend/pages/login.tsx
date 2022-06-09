import { NextPage } from "next";
import { useRouter } from "next/router";
import { useEffect } from "react";
import AppleLoginButton from "../components/Login/SocialLoginButton/AppleLoginButton/AppleLoginButton";
import GoogleLoginButton from "../components/Login/SocialLoginButton/GoogleLoginButton/GoogleLoginButton";
import { useAuth } from "../hooks/useAuth";
import { useAuthService } from "../hooks/useAuthService";
import styles from "../styles/Login.module.scss";

const Login: NextPage = () => {
  const { loading, user } = useAuth();
  const authService = useAuthService();
  const router = useRouter();

  useEffect(() => {
    if (loading) {
      return;
    }

    if (user) {
      router.push("/dashboard");
    }
  }, [user, loading]);

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

  if (loading) {
    return <p>Loading...</p>;
  }

  if (user) {
    return <p>Logging in...</p>;
  }

  return (
    <>
      <h1>Login</h1>
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
    </>
  );
};

export default Login;
