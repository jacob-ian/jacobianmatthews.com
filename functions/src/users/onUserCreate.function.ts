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
    createdAt: new Date(),
    updatedAt: new Date(),
  });
});

async function createUserDocument(user: User): Promise<void> {
  await firestore.doc(`users/${user.uid}`).create(user);
}

async function sendEmailVerificationLink(
  email: string,
  name: string,
): Promise<void> {
  const verificationLink = await admin
    .auth()
    .generateEmailVerificationLink(email);
  functions.logger.info(`Sending verification email to ${email}`);
  return emailService.sendEmail({
    to: email,
    message: {
      subject: "Verify Your Email Address | Jacob Ian Matthews",
      html: `Hey ${name},
      

      Please verify your email address by following the link <a href="${verificationLink}">${verificationLink}</a>.

      If this wasn't you, please ignore this email.


      Thanks,

      Jacob Ian Matthews
      <a href="https://jacobianmatthews.com">jacobianmatthews.com</a>
      `,
    },
  });
}
