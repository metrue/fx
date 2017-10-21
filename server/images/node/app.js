const Koa = require('koa');
const func = require('./function');

const app = new Koa();
app.use(ctx => {
    const msg = func(ctx.req);
    ctx.body = msg;
});

app.listen(3000);
