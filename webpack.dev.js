const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const webpack = require('webpack');

const HtmlWebPackPlugin = require("html-webpack-plugin");

const htmlPlugin = new HtmlWebPackPlugin({
    template: './src/index.html',
    filename: './index.html',
    favicon: 'src/favicon.png'
});

let plugins = [htmlPlugin];
if (process.env.PIXLET_BACKEND === "wasm") {
    plugins.push(
        new webpack.DefinePlugin({
            'PIXLET_WASM': JSON.stringify(true),
            'PIXLET_API_BASE': JSON.stringify('pixlet'),
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
    mode: 'development',
    devtool: 'source-map',
    devServer: {
        port: 3000,
        historyApiFallback: true,
        proxy: [
            {
                context: ['/api'],
                target: 'http://localhost:8080',
                ws: true,
            },
        ],
    },
    plugins: plugins,
});