import { HttpException, HttpStatus } from "./HttpException";

export class ForbiddenException extends HttpException {
  constructor(message?: string) {
    super(HttpStatus.FORBIDDEN, message || "Forbidden");
  }
}
