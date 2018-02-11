// This is required for Firefox Nightly, which ships with Custom Element
// support, but without Shadow DOM.
// https://github.com/webcomponents/webcomponentsjs/issues/874


/**
 * @license
 * Copyright (c) 2017 The Polymer Project Authors. All rights reserved.
 * This code may only be used under the BSD style license found at http://polymer.github.io/LICENSE.txt
 * The complete set of authors may be found at http://polymer.github.io/AUTHORS.txt
 * The complete set of contributors may be found at http://polymer.github.io/CONTRIBUTORS.txt
 * Code distributed by Google as part of the polymer project is also
 * subject to an additional IP rights grant found at http://polymer.github.io/PATENTS.txt
 */

(function () {
    /* eslint-disable */
    'use strict';
    // global for (1) existence means `WebComponentsReady` will file,
    // (2) WebComponents.ready == true means event has fired.
    window.WebComponents = window.WebComponents || {};
    // Feature detect which polyfill needs to be imported.
    var polyfills = [];
    if (!('import' in document.createElement('link'))) {
        polyfills.push('hi');
    }
    if (!('attachShadow' in Element.prototype && 'getRootNode' in Element.prototype) ||
        (window.ShadyDOM && window.ShadyDOM.force)) {
        polyfills.push('sd');
    }
    if (!window.customElements || window.customElements.forcePolyfill || polyfills.indexOf('sd') > -1) {
        polyfills.push('ce');
    }
    // NOTE: any browser that does not have template or ES6 features
    // must load the full suite (called `lite` for legacy reasons) of polyfills.
    if (!('content' in document.createElement('template')) || !window.Promise || !Array.from ||
        // Edge has broken fragment cloning which means you cannot clone template.content
        !(document.createDocumentFragment().cloneNode() instanceof DocumentFragment)) {
        polyfills = ['lite'];
    }

    if (polyfills.length) {
        var newScript = document.createElement('script');
        newScript.src = 'bower_components/webcomponentsjs/webcomponents-' + polyfills.join('-') + '.js';
        // NOTE: this is required to ensure the polyfills are loaded before
        // *native* html imports load on older Chrome versions. This *is* CSP
        // compliant since CSP rules must have allowed this script to run.
        // In all other cases, this can be async.
        if (document.readyState === 'loading' && ('import' in document.createElement('link'))) {
            document.write(newScript.outerHTML);
        } else {
            document.head.appendChild(newScript);
        }
    } else {
        // Ensure `WebComponentsReady` is fired also when there are no polyfills loaded.
        // however, we have to wait for the document to be in 'interactive' state,
        // otherwise a rAF may fire before scripts in <body>

        var fire = function () {
            requestAnimationFrame(function () {
                window.WebComponents.ready = true;
                document.dispatchEvent(new CustomEvent('WebComponentsReady', { bubbles: true }));
            });
        };

        if (document.readyState !== 'loading') {
            fire();
        } else {
            document.addEventListener('readystatechange', function wait() {
                fire();
                document.removeEventListener('readystatechange', wait);
            });
        }
    }
})();