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
    .pipe(gulp.dest('app'));
});

gulp.task('copy', ['vendor'], function() {
  return gulp.src([
      'app/**/*',
      '!app/bower_components/**/demo',
      '!app/bower_components/**/demo/**',
      '!app/bower_components/**/tests',
      '!app/bower_components/**/tests/**',
      '!app/bower_components/**/test',
      '!app/bower_components/**/test/**',
      '!app/bower_components/**/examples',
      '!app/bower_components/**/examples/**',
      '!app/**/.*',
      '!app/**/bower.json',
      '!app/**/*.md'
    ])
    .pipe(debug({
      title: 'copy:'
    }))
    .pipe(gulp.dest(dist));
});


gulp.task('minifyHTML', ['copy'], function() {
  return gulp.src([
      'app/**/*.html',
      '!app/**/app-router/app-router.html',
      '!app/**/app-router/app-router.csp.html',
      '!app/bower_components/web-component-tester/data/index.html',
      '!app/**/demo/**',
      '!app/**/test/**',
      '!app/**/tests/**',
      '!app/**/examples/**'
    ])
    .pipe(debug({
      title: 'minifyHTML:'
    }))
    .pipe(htmlminify())
    .pipe(gulp.dest("dist"))
});

gulp.task('minifyJS', ['minifyHTML'], function() {
  return gulp.src([
      'app/**/*.js',
      '!app/**/*min.js',
      '!app/**/platinum-push-messaging/service-worker.js',
      '!app/bower_components/async/support/sync-package-managers.js',
      '!app/bower_components/sw-toolbox/demo/service-worker.js',
      '!app/**/sw-toolbox/sw-toolbox.js',
      '!app/**/demo/**',
      '!app/**/test/**',
      '!app/**/tests/**',
      '!app/**/examples/**'
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
      'app/**/*.css',
      '!app/**/demo/**',
      '!app/**/test/**',
      '!app/**/tests/**',
      '!app/**/examples/**'
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
