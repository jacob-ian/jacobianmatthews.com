import * as admin from "firebase-admin";
import { Firestore } from "firebase-admin/firestore";

interface Email {
  to: string;
  message: {
    subject: string;
    html: string;
    text?: string;
  };
}

/**
 * The email service for jacobianmatthews.com.
 * This service relies on the Firebase Trigger Email extension:
 * https://firebase.google.com/products/extensions/firebase-firestore-send-email
 */
export class EmailService {
  private _firestore: Firestore;

  constructor() {
    this._firestore = admin.firestore();
  }

  public async sendEmail(email: Email): Promise<void> {
    await this._firestore.collection("mail").add(email);
  }
}
