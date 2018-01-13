/* eslint-env node */

'use strict';

module.exports = {
    staticFileGlobs: [
        'index.html',
        'bower_components/webcomponentsjs/webcomponents-*.js',
        'custom-webcomponents-loader.js',
        'images/*',
        'favicon.ico'
    ],
    navigateFallback: 'index.html',
    navigateFallbackWhitelist: [/^\/(projects|posts|impress)\//]
};
