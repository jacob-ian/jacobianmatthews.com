import * as functions from "firebase-functions";
import { AdminUserService } from "../users/admin-user.service";

const adminUserService = new AdminUserService();

export const onAdminEmailDelete = functions.firestore
  .document("admins/{documentId}")
  .onDelete(async (snapshot) => {
    const { email } = snapshot.data();
    return adminUserService.removeUserAdmin(email);
  });
