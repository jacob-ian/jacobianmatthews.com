import { HttpException, HttpStatus } from "./HttpException";

export class UnauthenticatedException extends HttpException {
  constructor(message?: string) {
    super(HttpStatus.UNAUTHENTICATED, message || "Unauthenticated");
  }
}
