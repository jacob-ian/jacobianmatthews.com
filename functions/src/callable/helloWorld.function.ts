import * as functions from "firebase-functions";

export const helloWorld = functions.https.onCall(() => {
  functions.logger.log("Hello world!");
});
