
/**
 * Module dependencies.
 */

var typeis = require('type-is');
var json = require('./json');
var form = require('./form');
var text = require('./text');

var jsonTypes = ['json', 'application/*+json', 'application/csp-report'];
var formTypes = ['urlencoded'];
var textTypes = ['text'];

/**
 * Return a Promise which parses form and json requests
 * depending on the Content-Type.
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
  opts = opts || {};

  // json
  var jsonType = opts.jsonTypes || jsonTypes;
  if (typeis(req, jsonType)) return json(req, opts);

  // form
  var formType = opts.formTypes || formTypes;
  if (typeis(req, formType)) return form(req, opts);

  // text
  var textType = opts.textTypes || textTypes;
  if (typeis(req, textType)) return text(req, opts);

  // invalid
  var type = req.headers['content-type'] || '';
  var message = type ? 'Unsupported content-type: ' + type : 'Missing content-type';
  var err = new Error(message);
  err.status = 415;
  return Promise.reject(err);
};
