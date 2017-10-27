
5.1.1 / 2017-03-24
==================

  * fix: getOptions change to clone
  * fix: ensure options are independent in each request

5.1.0 / 2017-03-21
==================

  * feat: add options to support return raw body (#56)

5.0.3 / 2017-03-19
==================

  * fix: ensure inflate in promise chain (#54)

5.0.2 / 2017-03-10
==================

  * fix: keep compatibility with qs@4 (#53)

5.0.1 / 2017-03-06
==================

  * dpes: qs@6.4.0

5.0.0 / 2017-03-02
==================

  * deps: upgrade qs to 6.x (#52)

4.2.0 / 2016-05-05
==================

  * test: test on node 4, 5, 6
  * feat: Added support for request body inflation

4.1.0 / 2016-05-05
==================

  * feat: form parse support custom qs module

4.0.0 / 2015-08-15
==================

  * Switch to Promises instead of thunks

3.1.0 / 2015-08-06
==================

 * travis: add v2, v3, remove 0.11
 * add custom types options
 * use type-is

3.0.0 / 2015-07-25
==================

 * Updated dependencies. Added qs options support via queryString option key. (@yanickrochon)
   * upgrade qs@4.0.0, raw-body@2.1.2

2.0.0 / 2015-05-04
==================

  * json parser support strict mode

1.2.0 / 2015-04-29
==================

 * Add JSON-LD as known JSON-Type (@vanthome)

1.1.0 / 2015-02-27
==================

 * Fix content-length zero should not parse json
 * Bump deps, qs@~2.3.3, raw-body@~1.3.3
 * add support for `text/plain`
 * json support for `application/json-patch+json`, `application/vnd.api+json` and `application/csp-report`
