import * as admin from "firebase-admin";
import * as functions from "firebase-functions";
import { Auth } from "firebase-admin/lib/auth/auth";
import { Firestore } from "firebase-admin/lib/firestore";
import { UserRecord } from "firebase-functions/v1/auth";
import { EmailService } from "../emails/email.service";

export const ADMIN_CLAIMS = { admin: true };

export class AdminUserService {
  private _auth: Auth;
  private _firestore: Firestore;
  private _emailService: EmailService;

  constructor() {
    this._auth = admin.auth();
    this._firestore = admin.firestore();
    this._emailService = new EmailService();
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
    return this._sendNewAdminEmail(email, user.displayName);
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

  private async _sendNewAdminEmail(
    email: string,
    name?: string,
  ): Promise<void> {
    return this._emailService.sendEmail({
      to: email,
      message: {
        subject: "Welcome, Admin! | Jacob Ian Matthews",
        text: `Hello${
          name ? " " + name : ""
        }, You have been made an Admin at jacobianmatthews.com. Start Editing...`,
        html: `
          Hello${name ? " " + name : ""},


          You have been made an Admin at <a href="https://jacobianmatthews.com>jacobianmatthews.com</a>.

          Start editing here:

          <a href="https://jacobianmatthews.com/dashboard">jacobianmatthews.com/dashboard</a>


          All the best,

          Jacob Ian Matthews
          <a href="https://jacobianmatthews.com>jacobianmatthews.com</a>
        `,
      },
    });
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
