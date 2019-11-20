;(function () {
    'use strict';

    function addClass(element, value) {
        element.classList.add(value);
    }

    function click(element, callback) {
        element.addEventListener('click', callback);
    }

    function clipboard(input) {
        var element = document.createElement('textarea');
        var previouslyFocusedElement = document.activeElement;
        element.value = input;
        // Prevent keyboard from showing on mobile
        element.setAttribute('readonly', '');
        element.style.contain = 'strict';
        element.style.position = 'absolute';
        element.style.left = '-9999px';
        element.style.fontSize = '12pt'; // Prevent zooming on iOS
        var selection = document.getSelection();
        var originalRange = false;
        if (selection.rangeCount > 0) {
            originalRange = selection.getRangeAt(0);
        }
        document.body.append(element);
        element.select();
        // Explicit selection workaround for iOS
        element.selectionStart = 0;
        element.selectionEnd = input.length;
        var isSuccess = false;
        try {
            isSuccess = document.execCommand('copy');
        } catch (_) {
        }
        element.remove();
        if (originalRange) {
            selection.removeAllRanges();
            selection.addRange(originalRange);
        }
        // Get the focus back on the previously focused element, if any
        if (previouslyFocusedElement) {
            previouslyFocusedElement.focus();
        }
        return isSuccess;
    }

    function delStyle(element) {
        element.removeAttribute('style');
    }

    function display(element, value) {
        if (isUndefined(value)) value = 'block';
        element.style.display = value;
    }

    function dumpSize(element) {
        console.log(
            // '\nelement.accessKey,' + element.accessKey +
            // '\nelement.animate,' + element.animate +
            // '\nelement.attachShadow,' + element.attachShadow +
            // '\nelement.attributes,' + element.attributes +
            // '\nelement.classList,' + element.classList +
            // '\nelement.className,' + element.className +
            '\nelement.clientHeight,' + element.clientHeight +
            '\nelement.clientLeft,' + element.clientLeft +
            '\nelement.clientTop,' + element.clientTop +
            '\nelement.clientWidth,' + element.clientWidth +
            // '\nelement.closest,' + element.closest +
            // '\nelement.computedStyleMap,' + element.computedStyleMap +
            // '\nelement.getAttribute,' + element.getAttribute +
            // '\nelement.getAttributeNames,' + element.getAttributeNames +
            // '\nelement.getAttributeNode,' + element.getAttributeNode +
            '\nelement.getBoundingClientRect,' + element.getBoundingClientRect +
            '\nelement.getClientRects,' + JSON.stringify(element.getClientRects()) +
            // '\nelement.getElementsByClassName,' + element.getElementsByClassName +
            // '\nelement.getElementsByTagName,' + element.getElementsByTagName +
            // '\nelement.hasAttribute,' + element.hasAttribute +
            // '\nelement.hasAttributes,' + element.hasAttributes +
            // '\nelement.hasPointerCapture,' + element.hasPointerCapture +
            // '\nelement.id,' + element.id +
            // '\nelement.innerHTML,' + element.innerHTML +
            // '\nelement.insertAdjacentElement,' + element.insertAdjacentElement +
            // '\nelement.insertAdjacentHTML,' + element.insertAdjacentHTML +
            // '\nelement.insertAdjacentText,' + element.insertAdjacentText +
            // '\nelement.localName,' + element.localName +
            // '\nelement.matches,' + element.matches +
            // '\nelement.name,' + element.name +
            // '\nelement.namespaceURI,' + element.namespaceURI +
            '\nelement.namespaceURI,' + element.offsetTop +
            // '\nelement.onfullscreenchange,' + element.onfullscreenchange +
            // '\nelement.onfullscreenerror,' + element.onfullscreenerror +
            // '\nelement.openOrClosedShadowRoot,' + element.openOrClosedShadowRoot +
            // '\nelement.outerHTML,' + element.outerHTML +
            // '\nelement.prefix,' + element.prefix +
            // '\nelement.querySelector,' + element.querySelector +
            // '\nelement.querySelectorAll,' + element.querySelectorAll +
            // '\nelement.releasePointerCapture,' + element.releasePointerCapture +
            // '\nelement.removeAttribute,' + element.removeAttribute +
            // '\nelement.removeAttributeNode,' + element.removeAttributeNode +
            // '\nelement.requestFullscreen,' + element.requestFullscreen +
            // '\nelement.requestPointerLock,' + element.requestPointerLock +
            // '\nelement.scroll,' + element.scroll +
            // '\nelement.scrollBy,' + element.scrollBy +
            '\nelement.scrollHeight,' + element.scrollHeight +
            // '\nelement.scrollIntoView,' + element.scrollIntoView +
            '\nelement.scrollLeft,' + element.scrollLeft +
            // '\nelement.scrollTo,' + element.scrollTo +
            '\nelement.scrollTop,' + element.scrollTop +
            '\nelement.scrollWidth,' + element.scrollWidth +
            '\nelement.offsetHeight,' + element.offsetHeight +
            '\nelement.offsetLeft,' + element.offsetLeft +
            '\nelement.offsetTop,' + element.offsetTop +
            '\nelement.offsetWidth,' + element.offsetWidth
        );
        // '\nelement.setAttribute,' + element.setAttribute +
        // '\nelement.setAttributeNode,' + element.setAttributeNode +
        // '\nelement.setCapture,' + element.setCapture +
        // '\nelement.setPointerCapture,' + element.setPointerCapture +
        // '\nelement.shadowRoot,' + element.shadowRoot +
        // '\nelement.slot,' + element.slot +
        // '\nelement.tagName,' + element.tagName +
        // '\nelement.toggleAttribute,' + element.toggleAttribute);
    }

    function dumpWindow() {
        console.log(
            // '\nwindow.alert,' + window.alert +
            // '\nwindow.applicationCache,' + window.applicationCache +
            // '\nwindow.blur,' + window.blur +
            // '\nwindow.cancelAnimationFrame,' + window.cancelAnimationFrame +
            // '\nwindow.cancelIdleCallback,' + window.cancelIdleCallback +
            // '\nwindow.clearImmediate,' + window.clearImmediate +
            // '\nwindow.close,' + window.close +
            // '\nwindow.closed,' + window.closed +
            // '\nwindow.confirm,' + window.confirm +
            // '\nwindow.console,' + window.console +
            // '\nwindow.crypto,' + window.crypto +
            // '\nwindow.customElements,' + window.customElements +
            // '\nwindow.defaultStatus,' + window.defaultStatus +
            '\nwindow.devicePixelRatio,' + window.devicePixelRatio +
            // '\nwindow.directories,' + window.directories +
            // '\nwindow.document,' + window.document +
            // '\nwindow.event,' + window.event +
            // '\nwindow.focus,' + window.focus +
            // '\nwindow.frameElement,' + window.frameElement +
            // '\nwindow.frames,' + window.frames +
            // '\nwindow.fullScreen,' + window.fullScreen +
            // '\nwindow.getComputedStyle,' + window.getComputedStyle +
            // '\nwindow.getSelection,' + window.getSelection +
            // '\nwindow.history,' + window.history +
            '\nwindow.innerHeight,' + window.innerHeight +
            '\nwindow.innerWidth,' + window.innerWidth +
            // '\nwindow.isSecureContext,' + window.isSecureContext +
            // '\nwindow.length,' + window.length +
            // '\nwindow.localStorage,' + window.localStorage +
            // '\nwindow.location,' + window.location +
            // '\nwindow.matchMedia,' + window.matchMedia +
            // '\nwindow.minimize,' + window.minimize +
            // '\nwindow.moveBy,' + window.moveBy +
            // '\nwindow.moveTo,' + window.moveTo +
            // '\nwindow.name,' + window.name +
            // '\nwindow.onappinstalled,' + window.onappinstalled +
            // '\nwindow.onbeforeinstallprompt,' + window.onbeforeinstallprompt +
            // '\nwindow.ondevicelight,' + window.ondevicelight +
            // '\nwindow.ondevicemotion,' + window.ondevicemotion +
            // '\nwindow.ondeviceorientation,' + window.ondeviceorientation +
            // '\nwindow.ondeviceproximity,' + window.ondeviceproximity +
            // '\nwindow.ondragdrop,' + window.ondragdrop +
            // '\nwindow.onuserproximity,' + window.onuserproximity +
            // '\nwindow.onvrdisplayactivate,' + window.onvrdisplayactivate +
            // '\nwindow.onvrdisplayblur,' + window.onvrdisplayblur +
            // '\nwindow.onvrdisplayconnect,' + window.onvrdisplayconnect +
            // '\nwindow.onvrdisplaydeactivate,' + window.onvrdisplaydeactivate +
            // '\nwindow.onvrdisplaydisconnect,' + window.onvrdisplaydisconnect +
            // '\nwindow.onvrdisplayfocus,' + window.onvrdisplayfocus +
            // '\nwindow.onvrdisplaypresentchange,' + window.onvrdisplaypresentchange +
            // '\nwindow.open,' + window.open +
            // '\nwindow.opener,' + window.opener +
            '\nwindow.outerHeight,' + window.outerHeight +
            '\nwindow.outerWidth,' + window.outerWidth +
            '\nwindow.pageXOffset,' + window.pageXOffset +
            '\nwindow.pageYOffset,' + window.pageYOffset +
            // '\nwindow.parent,' + window.parent +
            // '\nwindow.performance,' + window.performance +
            // '\nwindow.postMessage,' + window.postMessage +
            // '\nwindow.print,' + window.print +
            // '\nwindow.prompt,' + window.prompt +
            // '\nwindow.requestAnimationFrame,' + window.requestAnimationFrame +
            // '\nwindow.requestIdleCallback,' + window.requestIdleCallback +
            // '\nwindow.resizeBy,' + window.resizeBy +
            // '\nwindow.resizeTo,' + window.resizeTo +
            // '\nwindow.restore,' + window.restore +
            // '\nwindow.routeEvent,' + window.routeEvent +
            // '\nwindow.screen,' + window.screen +
            '\nwindow.screenLeft,' + window.screenLeft +
            '\nwindow.screenTop,' + window.screenTop +
            '\nwindow.screenX,' + window.screenX +
            '\nwindow.screenY,' + window.screenY
            // '\nwindow.scroll,' + window.scroll +
            // '\nwindow.scrollBy,' + window.scrollBy +
            // '\nwindow.scrollTo,' + window.scrollTo +
            // '\nwindow.scrollX,' + window.scrollX +
            // '\nwindow.scrollY,' + window.scrollY +
            // '\nwindow.self,' + window.self +
            // '\nwindow.sessionStorage,' + window.sessionStorage +
            // '\nwindow.showModalDialog,' + window.showModalDialog +
            // '\nwindow.speechSynthesis,' + window.speechSynthesis +
            // '\nwindow.status,' + window.status +
            // '\nwindow.stop,' + window.stop +
            // '\nwindow.top,' + window.top
            // '\nwindow.visualViewport,' + window.visualViewport +
            // '\nwindow.window,' + window.window
        )
    }

    function every(parent, selector, callback) {
        if (isUndefined(callback) && isFunction(selector)) {
            callback = selector;
        }
        if (isString(parent)) {
            selector = parent;
            parent = document;
        }
        return Array.prototype.slice.call(parent.querySelectorAll(selector)).map(function (element) {
            return callback(element);
        });
    }

    function interceptClick(element, callback) {
        element.addEventListener('click', function (event) {
            event.stopPropagation();
            callback && callback(event);
        })
    }

    function isFunction(value) {
        return typeof value === 'function';
    }

    function isString(value) {
        return typeof value === 'string';
    }

    function isUndefined(value) {
        return typeof value === 'undefined';
    }

    function isWeiXin() {
        return window.navigator.userAgent.toLowerCase().indexOf('micromessenger') !== -1;
    }

    function isWhitespace(text) {
        return !(isString(text) && text.trim().length > 0);
    }

    function removeClass(element, value) {
        element.classList.remove(value);
    }

    function substringAfter(text, delimiter) {
        var index = text.indexOf(delimiter);
        if (index === -1) return text;
        return text.substr(index + 1);
    }

    function substringAfterLast(s, c) {
        var i = s.lastIndexOf(c);
        if (i === -1) return s;
        return s.substring(i + c.length);
    }

    function substringBefore(text, delimiter) {
        var index = text.indexOf(delimiter);
        if (index === -1) return text;
        return text.substr(0, index);
    }

    function substringBeforeLast(s, c) {
        var i = s.lastIndexOf(c);
        if (i === -1) return s;
        return s.substring(0, i);
    }

    function toast(innerHTML, callback) {
        var mask = document.createElement('div');
        mask.setAttribute('style', 'position:fixed;top:0;bottom:0;left:0;right:0;background:rgba(0,0,0,.5);z-index:100');
        var container = document.createElement('div');
        container.setAttribute('style', 'width:100%;align-items:center;justify-content:center;height:100%;display:flex;');
        var messsage = document.createElement('div');
        messsage.setAttribute('style', 'width:50%;background:#fefefe;text-align:center;padding:1rem;');
        messsage.innerHTML = innerHTML;
        container.appendChild(messsage);
        mask.appendChild(container);
        mask.addEventListener('click', function () {
            mask.parentNode.removeChild(mask);
            callback && callback();
        });
        document.body.appendChild(mask);
    }

    function touchServer(options) {
        if (!options || !options.uri) return;
        var headers = options.headers || [];
        var method = options.method || 'GET';
        if (options.body) {
            if (!options.headers)
                headers.push(["Content-Type", "application/json"]);
            if (!options.method)
                method = 'POST';
        }
        console.log(options, method);

        fetch(options.uri, {
            method: method,
            headers: headers,
        }).then(function (response) {
            options.response && options.response(response);
            if (!response.ok) {
                throw new Error();
            }
            return response.json();
        }).then(function (data) {
            options.success && options.success(data);
        }).catch(function (error) {
            options.failed && options.failed(error);
        })
    }

    window['_'] = {
        addClass: addClass,
        click: click,
        clipboard: clipboard,
        delStyle: delStyle,
        display: display,
        dumpSize: dumpSize,
        dumpWindow: dumpWindow,
        every: every,
        interceptClick: interceptClick,
        isFunction: isFunction,
        isString: isString,
        isUndefined: isUndefined,
        isWeiXin: isWeiXin,
        isWhitespace: isWhitespace,
        removeClass: removeClass,
        substringAfter: substringAfter,
        substringAfterLast: substringAfterLast,
        substringBefore: substringBefore,
        substringBeforeLast: substringBeforeLast,
        toast: toast,
        touchServer: touchServer,
    }
})();