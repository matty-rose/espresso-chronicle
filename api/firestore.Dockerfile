FROM node:alpine

RUN apk add openjdk11

RUN npm install -g firebase-tools

WORKDIR /app

COPY firebase.json .

CMD [ "firebase", "--project=espressolog", "emulators:start", "--only", "firestore"]
