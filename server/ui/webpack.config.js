const path = require('path');

module.exports = {
    entry: './src/static/js/app.js',
    output: {
        path: path.resolve(__dirname, 'src/static/js'),
        filename: 'bundle.js'
    },
    mode: 'production',
};
