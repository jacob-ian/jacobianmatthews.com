import { AuthService } from "./AuthService";

describe("AuthService", () => {
  describe("handleAuthRedirect", () => {
    it.todo(
      "Should throw a NoAuthRedirectException if called without an auth redirect waiting",
    );

    it.todo("Should throw an error if client is missing csrfToken cookie");

    it.todo("Should call backend to log in with ID and csrf tokens");

    it.todo(
      "Should call firebase auth signOut after logging in to the backend",
    );
  });

  describe("signInWithGoogle", () => {
    it.todo("Should create a GoogleAuthProvider");

    it.todo("Should redirect the page to Google login page");
  });

  describe("signInWithApple", () => {
    it.todo("Should create an OAuthProvider with id apple.com");

    it.todo("Should redirect the page to Apple login page");
  });

  describe("signOut", () => {
    it.todo("Should call backend to logout");
  });
});
