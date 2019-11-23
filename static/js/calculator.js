;(function () {
    "use strict";

    var modContainer_ = document.querySelector('.mod-container');
    var modTabWrapper_ = document.querySelector('.mod-tab-wrapper');
    var modBtns_ = document.querySelectorAll('.mod-btn[data-value]');
    var modInput_ = document.querySelector('.mod-input');
    var clear_ = document.getElementById('clear');
    var save_ = document.getElementById('save');
    var modList_ = document.querySelector('.mod-list');
    var backspace = document.getElementById('backspace');


    var obj_;

    modContainer_.style.height = (window.innerHeight - modTabWrapper_.getClientRects()[0].height) + 'px';


    // ==============================================

    function setupButtons() {

        for (var i = 0; i < modBtns_.length; i++) {
            var btn = modBtns_[i];
            btn.addEventListener('click', function (event) {
                var value = event.currentTarget.getAttribute('data-value');
                var expression = substringBeforeLast(modInput_.textContent, '=').trim();
                expression += value;
                console.log(expression);
                try {
                    var results = eval(expression);
                    if (!isNaN(results)) {
                        modInput_.textContent = expression + ' = ' + results.toFixed(2);
                    } else {
                        modInput_.textContent = expression;
                    }
                } catch (e) {
                    modInput_.textContent = expression;
                }

            });
        }
    }

    function calculate() {
        var expression = formatExpression(modInput_.textContent);
        var sel = window.getSelection();
        var anchorOffset = sel.anchorOffset;
        expression = substringBeforeLast(expression, '=');
        var result = eval(expression);
        if (!isNaN(result)) {
            modInput_.textContent = expression.trim() + ' = ' + result.toFixed(2);
            var selection = window.getSelection();
            selection.collapse(modInput_.childNodes[modInput_.childNodes.length - 1], anchorOffset);
        }
    }


    function formatExpression(expression) {
        expression = expression.replace(/[^0-9.+-/=*]/g, '');
        expression = expression.replace('/[*]{2,}/g', '*');
        expression = expression.replace('/[+]{2,}/g', '*');
        expression = expression.replace('/[-]{2,}/g', '*');
        expression = expression.replace('/[/]{2,}/g', '*');
        for (var i = 0; i < expression.length; i++) {
        }
        return expression;
    }

    function initialize() {
        setupEditable();
        setupClear();
        setupSave();
        loadStorage();
        setupDeleteButtons();
        setupButtons();
        setBackSpace()
    }

    function setupClear() {
        clear_.addEventListener('click', function () {
            modInput_.textContent = '';
        });
    }

    function setupDeleteButtons() {
        var buttons = document.querySelectorAll('.mod-list-icon');
        for (var i = 0; i < buttons.length; i++) {
            buttons[i].onclick = function (event) {
                var t = event.currentTarget.parentNode.textContent.trim();
                var n = event.currentTarget.parentNode.parentNode;
                n.parentNode.removeChild(n);
                var index = obj_.indexOf(t);
                if (index !== -1) {
                    obj_.splice(index, 1);

                    if(_.isWeiXin()){
                        setCookie("calculator", JSON.stringify(obj_), 365);
                    }else {
                        window.localStorage && window.localStorage.setItem("calculator", JSON.stringify(obj_));
                    }

                }
            };
        }
    }

    function setupEditable() {
        modInput_.addEventListener('input', function () {
            calculate();
        });
    }

    function substringBeforeLast(s, c) {
        if (!s) return s;
        var i = s.lastIndexOf(c);
        if (i === -1) return s;
        return s.substring(0, i);
    }

    function onSave() {
        if (!obj_) obj_ = [];
        var value = modInput_.textContent;
        if (obj_.indexOf(value) === -1) {
            obj_.push(value);
            modList_.insertAdjacentHTML('beforeend', '<li class="mod-list-item">\n' +
                '            <a>\n' +
                '                <div class="mod-list-icon">\n' +
                '                    <svg height="16" viewBox="0 0 12 16" version="1.1" width="12" role="img">\n' +
                '                        <path fill-rule="evenodd" d="M7.48 8l3.75 3.75-1.48 1.48L6 9.48l-3.75 3.75-1.48-1.48L4.52 8 .77 4.25l1.48-1.48L6 6.52l3.75-3.75 1.48 1.48L7.48 8z"></path>\n' +
                '                    </svg>\n' +
                '                </div>\n' +
                value +
                '                </a>\n' +
                '        </li>')
        }
        if (_.isWeiXin()) {
            setCookie("calculator", JSON.stringify(obj_), 365);
        } else {
            window.localStorage && window.localStorage.setItem("calculator", JSON.stringify(obj_));
        }
        setupDeleteButtons();
    }


    function setupSave() {
        save_.addEventListener('click', onSave);
    }

    function loadStorage() {
        if (!_.isWeiXin() && !window.localStorage) return;
        var data = _.isWeiXin() ? getCookie('calculator') : window.localStorage.getItem('calculator');
        if (!data) return;
        obj_ = JSON.parse(data);
        for (var i = 0; i < obj_.length; i++) {
            modList_.insertAdjacentHTML('beforeend', '<li class="mod-list-item">\n' +
                '            <a>\n' +
                '                <div class="mod-list-icon">\n' +
                '                    <svg height="16" viewBox="0 0 12 16" version="1.1" width="12" role="img">\n' +
                '                        <path fill-rule="evenodd" d="M7.48 8l3.75 3.75-1.48 1.48L6 9.48l-3.75 3.75-1.48-1.48L4.52 8 .77 4.25l1.48-1.48L6 6.52l3.75-3.75 1.48 1.48L7.48 8z"></path>\n' +
                '                    </svg>\n' +
                '                </div>\n' +
                obj_[i] +
                '                </a>\n' +
                '        </li>')
        }
    }

    function setCookie(c_name, value, expiredays) {

        var exdate = new Date()

        exdate.setDate(exdate.getDate() + expiredays)

        document.cookie = c_name + "=" + escape(value) +

            ((expiredays == null) ? "" : ";expires=" + exdate.toGMTString())

    }


    function getCookie(c_name) {

        if (document.cookie.length > 0) {

            var c_start = document.cookie.indexOf(c_name + "=")

            if (c_start !== -1) {

                c_start = c_start + c_name.length + 1

                var c_end = document.cookie.indexOf(";", c_start)

                if (c_end === -1) c_end = document.cookie.length

                return unescape(document.cookie.substring(c_start, c_end))

            }

        }

        return ""

    }

    function setBackSpace() {
        backspace.addEventListener('click', function () {
            var text = modInput_.textContent;
            if (!text) return;
            text = substringBeforeLast(text, '=').trim();
            if (!text) return;
            text = text.substring(0, text.length - 1);
            modInput_.textContent = text;

            var expression = formatExpression(modInput_.textContent);
            var result = eval(expression);
            if (!isNaN(result)) {
                modInput_.textContent = expression.trim() + ' = ' + result.toFixed(2);
            }
            var selection = window.getSelection();
            selection.collapse(modInput_.childNodes[modInput_.childNodes.length - 1], substringBeforeLast(modInput_.value, '=').trim().length - 1);
        });
    }

    // ==============================================

    initialize();

})();