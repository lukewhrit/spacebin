FROM node:latest

RUN git clone https://github.com/spacebin-org/server.git spacebin 

WORKDIR spacebin

RUN yarn add sqlite3

EXPOSE 7777:7777

CMD ["yarn", "start"]
