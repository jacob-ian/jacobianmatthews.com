export enum HttpStatus {
  BAD_REQUEST = 400,
  UNAUTHENTICATED = 401,
  FORBIDDEN = 403,
  NOT_FOUND = 404,
  INTERNAL = 500,
}

export class HttpException extends Error {
  public httpStatus: HttpStatus | number;

  constructor(status: HttpStatus | number, message: string) {
    super(message);
    this.httpStatus = status;
  }
}
