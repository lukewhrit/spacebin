FROM node:14
WORKDIR /etc/spacebin-api
COPY package*.json ./
COPY . .
RUN yarn
EXPOSE 7777
CMD [ "yarn", "start" ]