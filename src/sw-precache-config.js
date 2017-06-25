/* eslint-env node */

'use strict';

module.exports = {
    staticFileGlobs: [
        'index.html',
        'bower_components/webcomponentsjs/webcomponents-*.js',
        'images/*',
        'favicon.ico'
    ],
    navigateFallback: 'index.html',
    navigateFallbackWhitelist: [/^\/(projects|posts|impress)\//]
};
