#!/usr/bin/env node
'use strict';
const meow = require('meow');
const getPort = require('get-port');

meow(`
	Example
	  $ get-port
	  51402
`);

getPort().then(console.log);
