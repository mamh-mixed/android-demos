'use strict';

var gulp = require('gulp');
var conf = require('./conf');
var $ = require('gulp-load-plugins')();

function deployToServer(baseDir, username, hostname, destination) {
  console.log('Deploy to ' + username + '@' + hostname + ':' + destination);
	return gulp.src(baseDir + '/**')
		.pipe($.rsync({
			root: baseDir,
			hostname: hostname,
			username: username,
			progress: true,
			destination: destination
		}));
}

gulp.task('deploy', ['build'], function() {
	let username = 'webapp';
	let hostname = 'dev.ipay.so';
	let destination = '/opt/quickpay/admin';
	return deployToServer(conf.paths.dist, username, hostname, destination);
});
gulp.task('deploy:test', ['build'], function() {
  let username = 'webapp';
	let hostname = 'dev.ipay.so';
	let destination = '/opt/quickpay/admin';
	return deployToServer(conf.paths.dist, server);
});
gulp.task('deploy:product', ['build'], function() {
  let username = 'webapp';
	let hostname = 'dev.ipay.so';
	let destination = '/opt/quickpay/admin';
	return deployToServer(conf.paths.dist, server);
});
