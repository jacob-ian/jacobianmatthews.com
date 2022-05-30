export interface User {
  uid: string;
  name: string;
  email: string;
  emailVerified: boolean;
  photoUrl?: string;
  createdAt: Date;
  updatedAt: Date;
}
