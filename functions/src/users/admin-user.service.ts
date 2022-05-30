import * as admin from "firebase-admin";
import * as functions from "firebase-functions";
import { Auth } from "firebase-admin/lib/auth/auth";
import { Firestore } from "firebase-admin/lib/firestore";
import { UserRecord } from "firebase-functions/v1/auth";

export const ADMIN_CLAIMS = { admin: true };

export class AdminUserService {
  private _auth: Auth;
  private _firestore: Firestore;

  constructor() {
    this._auth = admin.auth();
    this._firestore = admin.firestore();
  }

  public async makeUserAdmin(email: string): Promise<void> {
    return this._setUserCustomClaims(email, ADMIN_CLAIMS);
  }

  private async _setUserCustomClaims(
    email: string,
    claims: Record<string, any>,
  ): Promise<void> {
    const user = await this._getUserByEmail(email);
    if (!user) {
      return;
    }
    await this._auth.setCustomUserClaims(user.uid, claims);
    await this._auth.revokeRefreshTokens(user.uid);
  }

  private async _getUserByEmail(
    email: string,
  ): Promise<UserRecord | undefined> {
    try {
      return await this._auth.getUserByEmail(email);
    } catch {
      functions.logger.info("User with admin email does not exist.");
      return undefined;
    }
  }

  public async removeUserAdmin(email: string): Promise<void> {
    const customClaims = { admin: null };
    return this._setUserCustomClaims(email, customClaims);
  }

  public async isAdminEmail(email: string): Promise<boolean> {
    const adminEmailListing = await this._firestore
      .collection("admins")
      .where("email", "==", email)
      .get();
    return adminEmailListing.docs.length === 1;
  }
}
