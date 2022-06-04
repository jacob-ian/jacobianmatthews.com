import { HttpException, HttpStatus } from "./HttpException";

export class InternalErrorException extends HttpException {
  constructor(message?: string) {
    super(HttpStatus.INTERNAL, message || "Server Error");
  }
}
