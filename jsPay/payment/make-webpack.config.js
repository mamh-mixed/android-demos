/*
 * @Author: dmyang
 * @Date:   2015-08-02 14:16:41
 * @Last Modified by:   dmyang
 * @Last Modified time: 2015-09-14 15:08:13
 */

'use strict';

// @see http://christianalfoni.github.io/javascript/2014/12/13/did-you-know-webpack-and-react-is-awesome.html
// @see https://github.com/webpack/react-starter/blob/master/make-webpack-config.js

var path = require('path');
var fs = require('fs');
var yargs = require('yargs').argv;

var webpack = require('webpack');
var _ = require('lodash');

var ExtractTextPlugin = require('extract-text-webpack-plugin');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var OpenBrowserPlugin = require('open-browser-webpack-plugin');
var StringReplacePlugin = require("string-replace-webpack-plugin");

var UglifyJsPlugin = webpack.optimize.UglifyJsPlugin;
var CommonsChunkPlugin = webpack.optimize.CommonsChunkPlugin;

var srcDir = path.resolve(process.cwd(), 'src');
var build = 'assets/';
var sourceMap = require('./src/sourcemap.json');

var excludeFromStats = [
	/node_modules[\\\/]/
];

var config = require('config'),
	gatwayServer = '';

if (process.env.NODE_ENV === 'production') {
	gatwayServer = config['PRODUCT_REMOTE_URL'];
} else {
	gatwayServer = config['TEST_REMOTE_URL'];
}

function makeConf(options) {
	options = options || {};

	var urlMap = require('./config/default.json');

	var debug = options.debug !== undefined ? options.debug : true;
	var entries = genEntries();
	var chunks = Object.keys(entries);

	var serverUrl = '';
	if (options.product) {
		serverUrl = urlMap['PRODUCT_REMOTE_URL'];
	} else {
		serverUrl = urlMap['TEST_REMOTE_URL'];
	}
	var config = {
		entry: entries,

		output: {
			// 在debug模式下，__build目录是虚拟的，webpack的dev server存储在内存里
			path: path.resolve(debug ? '__build' : build),
			filename: debug ? '[name].js' : 'js/[chunkhash:8].[name].min.js',
			chunkFilename: debug ? '[chunkhash:8].chunk.js' : 'js/[chunkhash:8].chunk.min.js',
			hotUpdateChunkFilename: debug ? '[id].[chunkhash:8].js' : 'js/[id].[chunkhash:8].min.js',
			publicPath: debug ? '/__build/' : ''
		},

		resolve: {
			root: [srcDir, './node_modules'],
			// alias: sourceMap,
			extensions: ['', '.js', '.css', '.scss', '.tpl', '.png', '.jpg']
		},

		resolveLoader: {
			root: path.join(__dirname, 'node_modules')
		},

		module: {
			// noParse: ['zepto'],
			loaders: [{
				test: /util\.js$/,
				loader: StringReplacePlugin.replace({
					replacements: [{
						pattern: /IN_GOD_WE_TRUST/,
						replacement: function(match, p1, offset, string) {
							// console.log('******', match, '1', p1, '2', offset, '3', string);
							return serverUrl;
						}
					}]
				})
			}, {
				test: /\.(png|jpg)$/,
				loader: 'url-loader?mimetype=image/png&limit=8192'
			}, {
				test: /\.handlebars$/,
				loader: "handlebars-template-loader"
			}, {
				test: /\.(woff|eot|ttf)$/i,
				loader: 'url?limit=10000&name=fonts/[hash:8].[name].[ext]'
			}, {
				test: /\.(tpl|ejs)$/,
				loader: 'ejs'
			}, {
				test: /\.js$/,
				exclude: /(node_modules|lib|util\.js)/,
				loader: 'babel',
				query: {
					presets: ['es2015']
				}
			}]
		},

		plugins: [
			// new CommonsChunkPlugin({
			// 	name: 'vendors',
			// 	chunks: chunks,
			// 	// Modules must be shared between all entries
			// 	minChunks: chunks.length // 提取所有chunks共同依赖的模块
			// }),
			new OpenBrowserPlugin({
				url: 'http://192.168.199.193:8080/webpack-dev-server/pay.html'
			}),
			new webpack.DefinePlugin({
				__GATEWATSERVER__: gatwayServer
			}),
			new StringReplacePlugin()
			// new CommonsChunkPlugin({
			// 	name: 'common-bc',
			// 	chunks: ['vendors', 'b', 'c'],
			// 	minChunks: 2
			// })
		],

		devServer: {
			historyApiFallback: true,
			hot: true,
			inline: true,
			progress: true,
			contentBase: build,
			port: 8000
		},
	};

	if (options.product) {
		// 发布到生产上的代码要压缩
		config.plugins.push(
			new webpack.optimize.UglifyJsPlugin()
		);

		config.devtool = 'source-map';
	}

	if (debug) {
		// 开发阶段，css直接内嵌
		var cssLoader = {
			test: /\.css$/,
			loader: 'style!css'
		};
		var sassLoader = {
			test: /\.scss$/,
			loader: 'style!css!sass'
		};

		config.module.loaders.push(cssLoader);
		config.module.loaders.push(sassLoader);
	} else {
		// 编译阶段，css分离出来单独引入
		var cssLoader = {
			test: /\.css$/,
			loader: ExtractTextPlugin.extract('style', 'css?minimize') // enable minimize
		};
		var sassLoader = {
			test: /\.scss$/,
			loader: ExtractTextPlugin.extract('style', 'css?minimize', 'sass')
		};

		config.module.loaders.push(cssLoader);
		config.module.loaders.push(sassLoader);
		config.plugins.push(
			new ExtractTextPlugin('css/[contenthash:8].[name].min.css', {
				// 当allChunks指定为false时，css loader必须指定怎么处理
				// additional chunk所依赖的css，即指定`ExtractTextPlugin.extract()`
				// 第一个参数`notExtractLoader`，一般是使用style-loader
				// @see https://github.com/webpack/extract-text-webpack-plugin
				allChunks: false
			})
		);

		// 自动生成入口文件，入口js名必须和入口文件名相同
		// 例如，a页的入口文件是a.html，那么在js目录下必须有一个a.js作为入口文件
		var pages = fs.readdirSync(srcDir);

		pages.forEach(function(filename) {
			var m = filename.match(/(.+)\.html$/);

			if (m) {
				// @see https://github.com/kangax/html-minifier
				var conf = {
					template: path.resolve(srcDir, filename),
					// @see https://github.com/kangax/html-minifier
					// minify: {
					//     collapseWhitespace: true,
					//     removeComments: true
					// },
					filename: filename
				};

				if (m[1] in config.entry) {
					conf.inject = 'body';
					conf.chunks = ['vendors', m[1]];
				}

				config.plugins.push(new HtmlWebpackPlugin(conf));
			}
		});

		// config.plugins.push(new UglifyJsPlugin());
	}

	return config;
}

function genEntries() {
	var jsDir = path.resolve(srcDir, 'js');
	var names = fs.readdirSync(jsDir);
	var map = {};

	names.forEach(function(name) {
		// console.log(name);
		var m = name.match(/(.+)\.js$/);
		var entry = m ? m[1] : '';
		var entryPath = entry ? path.resolve(jsDir, name) : '';

		if (entry) map[entry] = entryPath;
	});

	// console.log(map);
	return map;
}

module.exports = makeConf;
