import { HttpException, HttpStatus } from "./HttpException";

export class NotFoundException extends HttpException {
  constructor(message?: string) {
    super(HttpStatus.NOT_FOUND, message || "Not Found");
  }
}
