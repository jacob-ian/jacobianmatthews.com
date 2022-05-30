import * as functions from "firebase-functions";
import { AdminUserService } from "./admin-user.service";

const adminUserService = new AdminUserService();

export const onAdminEmailAdd = functions.firestore
  .document("admins/{documentId}")
  .onCreate(async (snapshot) => {
    const { email } = snapshot.data();
    return adminUserService.makeUserAdmin(email);
  });
