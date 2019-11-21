;(function () {
    "use strict";

    var modContainer_ = document.querySelector('.mod-container');
    var modTabWrapper_ = document.querySelector('.mod-tab-wrapper');
    var modBtns_ = document.querySelectorAll('.mod-btn[data-value]');
    var modInput_ = document.querySelector('.mod-input');
    var clear_ = document.getElementById('clear');
    var save_ = document.getElementById('save');
    var modList_ = document.querySelector('.mod-list');


    var obj_;

    modContainer_.style.height = (window.innerHeight - modTabWrapper_.getClientRects()[0].height) + 'px';


    // ==============================================

    function calculate() {
        var expression = formatExpression(modInput_.textContent);
        var sel = window.getSelection();
        var anchorOffset = sel.anchorOffset;
        expression = substringBeforeLast(expression, '=');
        console.log(expression);
        var result = eval(expression);
        if (!isNaN(result)) {
            modInput_.textContent = expression.trim() + ' = ' + result.toFixed(2);
            var range = document.createRange();
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
    }

    function setupClear() {
        clear_.addEventListener('click', function () {
            modInput_.textContent = '';
        });
    }

    function setupDeleteButtons() {
        var buttons = document.querySelectorAll('.mod-list-icon');
        for (var i = 0; i < buttons.length; i++) {
            buttons[i].addEventListener('click', function (event) {
                var t = event.currentTarget.parentNode.textContent.trim();
                var n = event.currentTarget.parentNode.parentNode;
                n.parentNode.removeChild(n);
                var index = obj_.indexOf(t);
                if (index !== -1) {
                    obj_.splice(index, 1);
                    window.localStorage && window.localStorage.setItem("calculator", JSON.stringify(obj_));

                }
            })
        }
    }

    function setupEditable() {
        modInput_.addEventListener('input', function () {
            calculate();
        });
    }

    function substringBeforeLast(s, c) {
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
        window.localStorage && window.localStorage.setItem("calculator", JSON.stringify(obj_));

    }


    function setupSave() {
        save_.addEventListener('click', onSave);
    }

    function loadStorage() {
        if (!window.localStorage) return;
        var data = window.localStorage.getItem('calculator');
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

    // ==============================================

    initialize();

})();