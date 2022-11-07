FROM node:16-alpine
WORKDIR /firebase

RUN apk update && apk add openjdk11
RUN npm install -g firebase-tools

COPY firebase.json .
COPY .firebaserc .

RUN firebase setup:emulators:ui

CMD ["firebase", "emulators:start"]
