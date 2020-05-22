FROM metrue/fx-node-base

COPY . .
EXPOSE 3000
RUN npm i swagger-stats express-to-koa
CMD ["node", "app.js"]
