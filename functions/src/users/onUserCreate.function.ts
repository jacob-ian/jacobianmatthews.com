import * as functions from "firebase-functions";
import * as admin from "firebase-admin";
import { User } from "./user.interface";
import { AdminUserService } from "./admin-user.service";
import { EmailService } from "../emails/email.service";

const adminUserService = new AdminUserService();
const emailService = new EmailService();
const firestore = admin.firestore();

export const onUserCreate = functions.auth.user().onCreate(async (user) => {
  const { uid, displayName, email, emailVerified, photoURL } = user;
  if (!email || !displayName) {
    return functions.logger.error(
      `Cannot create user ${uid}. Missing inputs.`,
      user,
    );
  }

  if (emailVerified && (await adminUserService.isAdminEmail(email))) {
    await adminUserService.makeUserAdmin(email);
  }

  if (!emailVerified) {
    await sendEmailVerificationLink(email, displayName);
  }

  return createUserDocument({
    uid,
    email,
    emailVerified,
    name: displayName,
    photoUrl: photoURL,
  });
});

async function createUserDocument(user: User): Promise<void> {
  await firestore.doc(`users/${user.uid}`).create(user);
}

async function sendEmailVerificationLink(
  email: string,
  name: string,
): Promise<void> {
  const verificationLink = admin.auth().generateEmailVerificationLink(email);
  return emailService.sendEmail({
    to: email,
    message: {
      subject: "Verify Your Email Address | Jacob Ian Matthews",
      html: `Hey ${name},
      

      Please verify your email address by clicking <a href="${verificationLink}">here</a>.

      You can also paste the following into your browser:

      <pre>${verificationLink}</pre>


      Thanks,

      Jacob Ian Matthews
      <a href="https://jacobianmatthews.com">jacobianmatthews.com</a>
      `,
    },
  });
}
