const say = require('./helper')

module.exports = (ctx) => {
  say("hi")
  ctx.body = 'hello world'
}
