import { initializeApp } from "firebase-admin";
import { exportFunctions } from "firebase-functions-exporter";

initializeApp();
module.exports = exportFunctions();
