# Base image
FROM node:18

WORKDIR /app/gateway

COPY package.json yarn.lock ./

RUN yarn install

COPY . .

RUN yarn build

EXPOSE 50052

CMD [ "node", "dist/index.js" ]
