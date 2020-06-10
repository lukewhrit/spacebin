FROM node:14
WORKDIR /opt/spacebin-api
COPY package.json yarn.lock ./
COPY . .
RUN yarn
EXPOSE 7777
CMD [ "yarn", "start" ]
