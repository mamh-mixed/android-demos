'use strict';

var gulp = require('gulp');
var yargs = require('yargs').argv;
var webpack = require('webpack');
var uglify = require('gulp-uglify');
var gulpIf = require('gulp-if');
var sourcemaps = require('gulp-sourcemaps');
var gutil = require('gulp-util');
var WebpackDevServer = require("webpack-dev-server");
var stream = require('webpack-stream');

var webpackConf = require('./webpack.config');
var webpackDevConf = require('./webpack-dev.config');
var webpackPrdConf = require('./webpack-prd.config');

var path = {
	HTML: 'src/index.html',
	ALL: ['src/**/*.html', 'src/**/*.js'],
	MINIFIED_OUT: 'build.min.js',
	DEST_SRC: 'dist/src',
	DEST_BUILD: 'dist/build',
	DEST: 'dist'
};

var src = process.cwd() + '/src';
var assets = process.cwd() + '/assets';

// js check
gulp.task('hint', function() {
	var jshint = require('gulp-jshint');
	var stylish = require('jshint-stylish');

	return gulp.src([
			'!' + src + '/js/lib/**/*.js',
			src + '/js/**/*.js'
		])
		.pipe(jshint())
		.pipe(jshint.reporter('default'));
	// .pipe(jshint.reporter(stylish));
});

// clean assets
gulp.task('clean', ['hint'], function() {
	var clean = require('gulp-clean');

	return gulp.src(assets, {
		read: true
	}).pipe(clean());
});

// run webpack pack
gulp.task('pack', ['clean'], function(done) {
	if (yargs.p) {
		// 生产打包
		webpack(webpackPrdConf, function(err, stats) {
			if (err) throw new gutil.PluginError('webpack', err);
			gutil.log('[webpack]', stats.toString({
				colors: true
			}));
			done();
		});
	} else {
		// 开发打包
		webpack(webpackConf, function(err, stats) {
			if (err) throw new gutil.PluginError('webpack', err);
			gutil.log('[webpack]', stats.toString({
				colors: true
			}));
			done();
		});
	}

});

// html process
gulp.task('default', ['pack'], function() {
	var replace = require('gulp-replace');
	var htmlmin = require('gulp-htmlmin');

	return gulp
		.src(assets + '/*.html')
		.pipe(replace(/<script(.+)?data-debug([^>]+)?><\/script>/g, ''))
		// @see https://github.com/kangax/html-minifier
		.pipe(htmlmin({
			collapseWhitespace: true,
			removeComments: true
		}))
		.pipe(gulp.dest(assets));
});

// deploy assets to remote server
// -p 发布到生产环境， 否则发布到测试环境
gulp.task('deploy', ['default'], function() {
	var rsync = require('gulp-rsync'),
		target = yargs.p ? 'product' : 'test';

	var deployMap = {
		'test': {
			'destination': '/home/weixin/cloudCashier/agent'
		},
		'product': {
			// 'destination': '/home/weixin/cloudCashier/agent'
			'destination': '/home/weixin/cloudCashier/payment'
		}
	};

	yargs.p ? console.log('**************DEPLOY TO PRODUCT:: ' + deployMap[target]['destination'] + '*****************') : console.log('**************DEPLOY TO TEST:: ' + deployMap[target]['destination'] + '*****************');

	return gulp.src(assets + '/**')
		.pipe(rsync({
			root: 'assets',
			hostname: '139.129.116.65',
			username: 'weixin',
			progress: true,
			destination: deployMap[target]['destination']
		}));
});

gulp.task('webpack', [], function() {
	return gulp.src(path.ALL)
		.pipe(sourcemaps.init())
		.pipe(stream(webpackConf))
		.pipe(uglify())
		.pipe(sourcemaps.write())
		.pipe(gulp.dest(path.DEST_BUILD));
});

gulp.task("webpack-dev-server", function(callback) {
	// modify some webpack config options
	var myConfig = Object.create(webpackConf);
	myConfig.devtool = "eval";
	myConfig.debug = true;

	// Start a webpack-dev-server
	new WebpackDevServer(webpack(myConfig), {
		publicPath: "/" + myConfig.output.publicPath,
		stats: {
			colors: true
		}
	}).listen(8080, "localhost", function(err) {
		if (err) throw new gutil.PluginError("webpack-dev-server", err);
		gutil.log("[webpack-dev-server]", "http://localhost:8080/webpack-dev-server/index.html");
	});
});

gulp.task('watch', function() {
	gulp.watch(path.ALL, ['webpack']);
});

gulp.task('dev', ['webpack-dev-server', 'watch']);
