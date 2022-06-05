import { FirebaseApp } from "firebase/app";
import {
  Auth,
  connectAuthEmulator,
  getAuth,
  getRedirectResult,
  GoogleAuthProvider,
  inMemoryPersistence,
  OAuthProvider,
  signInWithRedirect,
} from "firebase/auth";
import { User } from "../../interfaces/user.interface";
import { getCookie } from "../../utils/getCookie";
import { isDevEnvironment } from "../../utils/isDevEnvironment";
import { HttpService } from "../http/http.service";
import { InvalidAuthException } from "./InvalidAuthException";
import { NoAuthRedirectException } from "./NoAuthRedirectException";

export class AuthService {
  private _firebaseAuth: Auth;
  private _http: HttpService;

  constructor(firebase?: FirebaseApp) {
    this._http = new HttpService();
    this._firebaseAuth = getAuth(firebase);
    if (isDevEnvironment()) {
      connectAuthEmulator(
        this._firebaseAuth,
        process.env.FIREBASE_AUTH_EMULATOR_HOST || "http://localhost:9099",
        { disableWarnings: true },
      );
    }
    this._firebaseAuth.setPersistence(inMemoryPersistence);
  }

  public async signInWithGoogle(): Promise<void> {
    const googleProvider = new GoogleAuthProvider();
    return signInWithRedirect(this._firebaseAuth, googleProvider);
  }

  public async signInWithApple(): Promise<void> {
    const appleProvider = new OAuthProvider("apple.com");
    return signInWithRedirect(this._firebaseAuth, appleProvider);
  }

  public async handleAuthRedirect(): Promise<void> {
    const result = await getRedirectResult(this._firebaseAuth);
    if (!result) {
      throw new NoAuthRedirectException();
    }
    const idToken = await result.user.getIdToken();
    const csrfToken = getCookie("csrfToken");
    if (!csrfToken) {
      throw new InvalidAuthException("Missing CSRF Token");
    }
    await this._loginToBackendWithIdToken(idToken, csrfToken);
    await this._firebaseAuth.signOut();
  }

  private async _loginToBackendWithIdToken(
    idToken: string,
    csrfToken: string,
  ): Promise<void> {
    await this._http.post({
      url: "/api/auth/login",
      body: { idToken, csrfToken },
    });
  }

  public async getSignedInUser(): Promise<User> {
    return this._http.get<User>({
      url: "/api/auth/me",
    });
  }

  public async signOut(): Promise<void> {
    await this._http.post({
      url: "/api/auth/logout",
    });
  }
}
