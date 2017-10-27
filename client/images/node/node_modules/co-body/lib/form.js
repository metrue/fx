
/**
 * Module dependencies.
 */

var raw = require('raw-body');
var inflate = require('inflation');
var qs = require('qs');
var utils = require('./utils');

/**
 * Return a Promise which parses x-www-form-urlencoded requests.
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
  var queryString = opts.queryString || {};

  // keep compatibility with qs@4
  if (queryString.allowDots === undefined) queryString.allowDots = true;

  // defaults
  var len = req.headers['content-length'];
  var encoding = req.headers['content-encoding'] || 'identity';
  if (len && encoding === 'identity') opts.length = ~~len;
  opts.encoding = opts.encoding || 'utf8';
  opts.limit = opts.limit || '56kb';
  opts.qs = opts.qs || qs;

  // raw-body returns a Promise when no callback is specified
  return Promise.resolve()
    .then(function() {
      return raw(inflate(req), opts);
    })
    .then(function(str){
      try {
        var parsed = opts.qs.parse(str, queryString);
        return opts.returnRawBody ? { parsed: parsed, raw: str } : parsed;
      } catch (err) {
        err.status = 400;
        err.body = str;
        throw err;
      }
    });
};
