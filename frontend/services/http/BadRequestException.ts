import { HttpException, HttpStatus } from "./HttpException";

export class BadRequestException extends HttpException {
  constructor(message?: string) {
    super(HttpStatus.BAD_REQUEST, message || "Bad Request");
  }
}
