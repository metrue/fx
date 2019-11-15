const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const fx = require('./fx');

const app = new Koa();
app.use(bodyParser());
app.use(fx);

app.listen(3000);
