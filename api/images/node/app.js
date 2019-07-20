const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const func = require('./fx');

const app = new Koa();
app.use(bodyParser());
app.use(async (ctx) => {
  const msg = await func(ctx.request.body, ctx);
  ctx.body = msg;
});

app.listen(3000);
