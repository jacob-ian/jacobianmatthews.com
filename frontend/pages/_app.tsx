import "../styles/globals.css";
import type { AppProps } from "next/app";
import AuthProvider from "../context/auth/AuthProvider";

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <AuthProvider>
      <Component {...pageProps} />
    </AuthProvider>
  );
}
