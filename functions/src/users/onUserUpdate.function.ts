import * as functions from "firebase-functions";
import { AdminUserService } from "./admin-user.service";

const adminUserService = new AdminUserService();

export const onUserUpdate = functions.firestore
  .document("user/${userId}")
  .onUpdate(async (change) => {
    const { emailVerified, email } = change.after.data();
    if (emailVerified === false) {
      return;
    }
    if (await adminUserService.isAdminEmail(email)) {
      return adminUserService.makeUserAdmin(email);
    }
  });
