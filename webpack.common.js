const webpack = require('webpack');

let plugins = [];

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

module.exports = {
    plugins,
    resolve: {
        extensions: ['*', '.js', '.jsx'],
    },
    experiments: {
        asyncWebAssembly: true,
        syncWebAssembly: true
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader'
                }
            },
            {
                test: /\.css$/,
                use: [
                    'style-loader',
                    {
                        loader: 'css-loader',
                        options: {
                            modules: true,
                        },
                    },
                ],
            },
            {
                test: /\.(webp|jpe?g|gif|star)$/i,
                use: [
                    {
                        loader: 'file-loader',
                    },
                ],
            },
            {
                test: /\.svg$/,
                use: ['@svgr/webpack'],
            },
        ]
    },
};