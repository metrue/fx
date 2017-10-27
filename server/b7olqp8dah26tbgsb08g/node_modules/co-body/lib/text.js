/**
 * Module dependencies.
 */

var raw = require('raw-body');
var inflate = require('inflation');
var utils = require('./utils');

/**
 * Return a Promise which parses text/plain requests.
 *
 * Pass a node request or an object with `.req`,
 * such as a koa Context.
 *
 * @param {Request} req
 * @param {Options} [opts]
 * @return {Function}
 * @api public
 */

module.exports = function(req, opts){
  req = req.req || req;
  opts = utils.clone(opts);

  // defaults
  var len = req.headers['content-length'];
  var encoding = req.headers['content-encoding'] || 'identity';
  if (len && encoding === 'identity') opts.length = ~~len;
  opts.encoding = opts.encoding || 'utf8';
  opts.limit = opts.limit || '1mb';

  // raw-body returns a Promise when no callback is specified
  return Promise.resolve()
    .then(function() {
      return raw(inflate(req), opts);
    })
    .then(str => {
      // ensure return the same format with json / form
      return opts.returnRawBody ? { parsed: str, raw: str } : str;
    });
};
