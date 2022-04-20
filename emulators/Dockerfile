# REQUIRES ROOT BUILD CONTEXT
FROM spine3/firebase-emulator
WORKDIR /firebase
COPY functions .
WORKDIR /firebase/functions
RUN npm run build
ENTRYPOINT [ "/runner.sh" ]