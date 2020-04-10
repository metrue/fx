The only thing `fx` does is to make a single function to be an HTTP service, `fx` takes two steps to finish the process,
* Wraps the function into a web server project and make it be the handler function of the HTTP request.
* Builds the bundled web server to be Docker image, then run it as a Docker container and expose the service port.

`bundler` is the component responsible for step 1: wrap a single function (and its dependencies) into a web server.  Take a Node web service as an example, the project looks like this,
```
helloworld
├── Dockerfile
├── app.js
└── fx.js
```

And the codes is pretty simple,

**app.js**
```javascript
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
```

**fx.js**
```javascript
module.exports = (ctx) => {
  ctx.body = 'hello world'
}
```

**Dockerfile**
```dockerfile
FROM metrue/fx-node-base

COPY . .
EXPOSE 3000
CMD ["node", "app.js"]
```

You can see that it's a web service built with `Koa`,  the request handler function defined in `fx.js` to response plain text `hello world`,  so to make a general JavaScript/Node function to be a web service, we only have to put it into `fx.js` above, then build and run it with Docker and that's it.

To support different programming languages, `fx` needs to implement different `bundlers`, the reasons are,
* The way to set up a web service is different in different languages
* The way to manage dependency is different in different languages.

So there will be (are) different `bundlers` in `fx`.
```
go-bundler: based on Gin
ruby-bundler: based on Sinatra
node-bundler: based on Koa
python-bundler: based on flask
...
```


