export class NoAuthRedirectException extends Error {
  constructor() {
    super("Redirect was not caused by user sign in");
  }
}
