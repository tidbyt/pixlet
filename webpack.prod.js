const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const webpack = require('webpack');
const path = require('path');
const HtmlWebPackPlugin = require("html-webpack-plugin");
const CopyWebpackPlugin = require('copy-webpack-plugin');

const htmlPlugin = new HtmlWebPackPlugin({
    template: './src/index.html',
    filename: '../index.html',
    favicon: 'src/favicon.png'
});

const copyPlugin = new CopyWebpackPlugin({
    patterns: [
        { from: "public", to: "../" },
    ],
});


let plugins = [htmlPlugin, copyPlugin];
if (process.env.PIXLET_BACKEND === "wasm") {
    plugins.push(
        new webpack.DefinePlugin({
            'PIXLET_WASM': JSON.stringify(true),
            'PIXLET_API_BASE': JSON.stringify('static/pixlet'),
        })
    );
} else {
    plugins.push(
        new webpack.DefinePlugin({
            'PIXLET_WASM': JSON.stringify(false),
            'PIXLET_API_BASE': JSON.stringify(''),
        })
    );
}

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
    plugins: plugins,
});