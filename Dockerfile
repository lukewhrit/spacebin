FROM node:12.16.3

RUN git clone https://github.com/spacebin-org/server.git spacebin

WORKDIR /opt/spacebin

RUN yarn add sqlite3

EXPOSE 7777

CMD ["yarn", "start"]
