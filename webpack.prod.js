const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const path = require('path');
const HtmlWebPackPlugin = require("html-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");

const htmlPlugin = new HtmlWebPackPlugin({
    template: './src/index.html',
    filename: '../index.html',
    favicon: 'src/favicon.png'
});

const copyPlugin = new CopyPlugin({
    patterns: [
        { from: "src/dist.go.txt", to: "../dist.go" },
    ],
});

module.exports = merge(common, {
    mode: 'production',
    devtool: 'source-map',
    output: {
        asyncChunks: true,
        publicPath: '/static/',
        path: path.resolve(__dirname, 'dist/static'),
        filename: '[name].[chunkhash].js',
        clean: true,
    },
    plugins: [htmlPlugin, copyPlugin]
});