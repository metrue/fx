const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const cors = require('@koa/cors');
const swStats = require('swagger-stats');
const e2k = require('express-to-koa');
const fx = require('./fx');

const app = new Koa();

app.use(e2k(swStats.getMiddleware({})));
app.use(cors({
  origin: '*',
}));
app.use(bodyParser());
app.use(fx);

app.listen(3000);
