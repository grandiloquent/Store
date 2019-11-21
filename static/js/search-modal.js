;(function () {
    "use strict";
    var SearchModal = function SearchModal() {
    };
    window['SearchModal'] = SearchModal;
    SearchModal.prototype.initialize = function (closeCallback) {
        this.element_ = document.getElementById('search-modal');
        if (!this.element_) return;
        this.closeCallback_ = closeCallback;
        this.complete_ = document.querySelector('.search-modal_autocomplete');
        this.input_ = this.element_.querySelector('.top_input-text');
        this.close_ = this.element_.querySelector('.top_clear');
        this.input_.addEventListener('input', this.onInput.bind(this));
        _.click(this.close_, this.onClose.bind(this));
        this.topBack_ = this.element_.querySelector('.top_back');
        _.click(this.topBack_, this.onHide.bind(this));
        this.searchModal_ = this.element_.querySelector('.search-modal');
        this.searchModalTop_ = this.element_.querySelector('.search-modal_top');
        this.searchModalRecent_ = this.element_.querySelector('.search-modal_recent');
        this.searchModalTop_.addEventListener('transitionend', function () {
            console.log('transitionend')
            // webkitAnimationEnd oanimationend oAnimationEnd msAnimationEnd animationend
            _.addClass(this.searchModalRecent_, 'active');
        }.bind(this));

        this.modalHotTags_ = this.element_.querySelector('.modal_hot-tags');
        this.searchModalClear_ = this.element_.querySelector('.search-modal_clear');
        this.searchModalHistorytags_ = this.element_.querySelector('.search-modal_historytags');


        this.offset_ = 0;
        this.setupRefresh();
        this.setupHistory();
        this.setupHistoryClear();
        this.onRefresh();
    };
    SearchModal.prototype.show = function () {
        // _.addClass(document.querySelector('html'), 'noscroll');
        _.addClass(this.searchModal_, 'active');
        setTimeout(function () {
            _.addClass(this.searchModalTop_, 'active');
        }.bind(this), 1000 / 60);
    };
    SearchModal.prototype.onClose = function () {
        this.input_.value = '';
        this.clearSearch();
    };
    SearchModal.prototype.onHide = function () {
        _.removeClass(this.searchModal_, 'active');
        _.removeClass(this.searchModalTop_, 'active');
        _.removeClass(this.searchModalRecent_, 'active');
        this.closeCallback_ && this.closeCallback_();
    };
    SearchModal.prototype.clearSearch = function () {
        _.delStyle(this.close_);
        _.delStyle(this.complete_);
    };
    SearchModal.prototype.onInput = function () {
        if (_.isWhitespace(this.input_.value)) {
            this.clearSearch();
        } else {
            _.display(this.complete_);
            _.display(this.close_, 'flex');
        }
    }
    SearchModal.prototype.setupRefresh = function () {
        this.modalHotRefresh_ = this.element_.querySelector('.modal_hot-refresh');
        this.modalHotRefresh_.addEventListener('click', this.onRefresh.bind(this));
    };
    SearchModal.prototype.onRefresh = function () {
        _.touchServer({
            uri: '/store/api/search?method=fetch&limit=6&sorttype=2&offset=' + this.offset_,
            success: this.onRefreshSuccess.bind(this),
            failed: this.onRefreshFailed
        })
    };
    SearchModal.prototype.onRefreshSuccess = function (obj) {
        if (!obj) return null;
        this.offset_ += 6;
        var buf = [];
        obj.forEach(function (item) {
            buf.push('<div class="hot-tag">' + item + '</div>');
        });
        this.modalHotTags_.innerHTML = buf.join('');
        this.bindHotTagEvent();
    };

    SearchModal.prototype.onRefreshFailed = function () {

    };

    SearchModal.prototype.bindHotTagEvent = function () {
        var that = this;
        _.every(this.modalHotTags_, '.hot-tag', function (element) {
            element.addEventListener('click', function (event) {
                that.addHistory(event.currentTarget.textContent)
                window.location = "/store/results?keyword=" + encodeURIComponent(event.currentTarget.textContent);

            })
        })
    };
    SearchModal.prototype.addHistory = function (keywords) {
        if (!window.localStorage) return;
        var obj = (JSON.parse(window.localStorage.getItem('search-history'))) || [];
        if (obj.indexOf(keywords) === -1) {
            obj.push(keywords);
        }
        window.localStorage.setItem('search-history', JSON.stringify(obj));
    };
    SearchModal.prototype.fetchHistory = function () {
        //"<div class=\"history-tag\">"+history+"</div>"
        if (!window.localStorage) return;
        var searchHistory = window.localStorage.getItem('search-history');
        if (!searchHistory) return;
        var obj = JSON.parse(searchHistory);
        var content = [];
        for (var i = 0; i < obj.length; i++) {
            var pattern =
                "<div class=\"history-tag\">" + obj[i] + "</div>";
            content.push(pattern);
        }
        if (!content.length) return;
        console.log(content.join(''));
        this.searchModalHistorytags_.innerHTML = content.join('');
        var that = this;
        _.every(this.searchModalHistorytags_, '.history-tag', function (element) {
            element.addEventListener('click', function (event) {
                that.input_.value = event.currentTarget.textContent;
            });
        });
    };
    SearchModal.prototype.setupHistory = function () {
        this.fetchHistory();
        this.bindHotTagEvent();
    };
    SearchModal.prototype.setupHistoryClear = function () {
        this.searchModalClear_.addEventListener('click', function () {
            window.localStorage && window.localStorage.removeItem('search-history');
            this.modalHotTags_.innerHTML = '';
        }.bind(this));
    };
    SearchModal.prototype.setupInput = function () {
        this.input_.addEventListener('input', function () {

        });
    }
})();
