export class InvalidAuthException extends Error {
  constructor(message: string) {
    super(message);
  }
}
