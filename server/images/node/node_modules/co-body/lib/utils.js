
/**
 * Module dependencies.
 */

exports.clone = function (opts) {
  var options = {};
  opts = opts || {};
  for (var key in opts) {
    options[key] = opts[key];
  }
  return options;
}
