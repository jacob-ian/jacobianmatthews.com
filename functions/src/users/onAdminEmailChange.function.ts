import * as functions from "firebase-functions";
import { AdminUserService } from "./admin-user.service";

const adminUserService = new AdminUserService();

export const onAdminEmailChange = functions.firestore
  .document("admins/{documentId}")
  .onUpdate(async (change) => {
    const { email: emailBefore } = change.before.data();
    const { email: emailAfter } = change.after.data();

    await adminUserService.removeUserAdmin(emailBefore);

    if (!emailAfter) {
      return;
    }

    await adminUserService.makeUserAdmin(emailAfter);
  });
