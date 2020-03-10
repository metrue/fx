const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const cors = require('@koa/cors');
const fx = require('./fx');

const app = new Koa();
app.use(cors({
  origin: '*',
}));
app.use(bodyParser());
app.use(fx);

app.listen(3000);
