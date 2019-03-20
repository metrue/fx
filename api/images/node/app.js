const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const func = require('./fx');

const app = new Koa();
app.use(bodyParser());
app.use(ctx => {
  const msg = func(ctx.request.body);
  ctx.body = msg;
});

app.listen(3000);
