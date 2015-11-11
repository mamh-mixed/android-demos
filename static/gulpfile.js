var gulp = require('gulp');
var debug = require('gulp-debug');
var htmlminify = require("gulp-html-minify");
var minify = require('gulp-minify');
var minifyCSS = require('gulp-minify-css')
var del = require('del');

var dist = 'dist';

gulp.task('clean', function() {
  return del([dist]);
});

gulp.task('vendor', ['clean'], function() {
  return gulp.src(
      'vendor/**/*'
    )
    .pipe(debug({
      title: 'vendor:'
    }))
    .pipe(gulp.dest('v0.0.1'));
});

gulp.task('copy', ['vendor'], function() {
  return gulp.src([
      'v0.0.1/**/*',
      '!v0.0.1/bower_components/**/demo',
      '!v0.0.1/bower_components/**/demo/**',
      '!v0.0.1/bower_components/**/tests',
      '!v0.0.1/bower_components/**/tests/**',
      '!v0.0.1/bower_components/**/test',
      '!v0.0.1/bower_components/**/test/**',
      '!v0.0.1/bower_components/**/examples',
      '!v0.0.1/bower_components/**/examples/**',
      '!v0.0.1/**/.*',
      '!v0.0.1/**/bower.json',
      '!v0.0.1/**/*.md'
    ])
    .pipe(debug({
      title: 'copy:'
    }))
    .pipe(gulp.dest(dist));
});


gulp.task('minifyHTML', ['copy'], function() {
  return gulp.src([
      'v0.0.1/**/*.html',
      '!v0.0.1/**/app-router/app-router.html',
      '!v0.0.1/**/app-router/app-router.csp.html',
      '!v0.0.1/bower_components/web-component-tester/data/index.html',
      '!v0.0.1/**/demo/**',
      '!v0.0.1/**/test/**',
      '!v0.0.1/**/tests/**',
      '!v0.0.1/**/examples/**'
    ])
    .pipe(debug({
      title: 'minifyHTML:'
    }))
    .pipe(htmlminify())
    .pipe(gulp.dest("dist"))
});

gulp.task('minifyJS', ['minifyHTML'], function() {
  return gulp.src([
      'v0.0.1/**/*.js',
      '!v0.0.1/**/*min.js',
      '!v0.0.1/**/platinum-push-messaging/service-worker.js',
      '!v0.0.1/bower_components/async/support/sync-package-managers.js',
      '!v0.0.1/bower_components/sw-toolbox/demo/service-worker.js',
      '!v0.0.1/**/sw-toolbox/sw-toolbox.js',
      '!v0.0.1/**/demo/**',
      '!v0.0.1/**/test/**',
      '!v0.0.1/**/tests/**',
      '!v0.0.1/**/examples/**'
    ])
    .pipe(debug({
      title: 'minifyJS:'
    }))
    .pipe(minify())
    .pipe(gulp.dest(dist))
});

gulp.task('minifyCSS', ['minifyHTML'], function() {
  // 1. 找到文件
  return gulp.src([
      'v0.0.1/**/*.css',
      '!v0.0.1/**/demo/**',
      '!v0.0.1/**/test/**',
      '!v0.0.1/**/tests/**',
      '!v0.0.1/**/examples/**'
    ])
    .pipe(debug({
      title: 'minifyCSS:'
    }))
    // 2. 压缩文件
    .pipe(minifyCSS())
    // 3. 另存为压缩文件
    .pipe(gulp.dest(dist))
});

gulp.task('default', ['clean', 'vendor', 'copy', 'minifyHTML', 'minifyJS', 'minifyCSS']);
