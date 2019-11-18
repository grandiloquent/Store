;(function () {
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
        console.log(this.searchModal_);
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
})();
/*
this.active_ = this.element_.querySelector('.active');
this.autocompleteListItem_ = this.element_.querySelector('.autocomplete-list_item');
this.close_ = this.element_.querySelector('.close');
this.copyRight_ = this.element_.querySelector('.copy-right');
this.flexCenter_ = this.element_.querySelector('.flex-center');
this.footer_ = this.element_.querySelector('.footer');
this.footerBar_ = this.element_.querySelector('.footer-bar');
this.footerBarWrapper_ = this.element_.querySelector('.footer-bar-wrapper');
this.header_ = this.element_.querySelector('.header');
this.hide_ = this.element_.querySelector('.hide');
this.historyTag_ = this.element_.querySelector('.history-tag');
this.homeRecommend_ = this.element_.querySelector('.home-recommend');
this.homeRecommendItem_ = this.element_.querySelector('.home-recommend-item');
this.homeRecommendItemMain_ = this.element_.querySelector('.home-recommend-item-main');
this.homeSearch_ = this.element_.querySelector('.home-search');
this.homeSearchFixed_ = this.element_.querySelector('.home-search-fixed');
this.homeSearchInputIcon_ = this.element_.querySelector('.home-search-input-icon');
this.homeSearchInputText_ = this.element_.querySelector('.home-search-input-text');
this.hotTag_ = this.element_.querySelector('.hot-tag');
this.icon_ = this.element_.querySelector('.icon');
this.likeCell_ = this.element_.querySelector('.like-cell');
this.likeCellBottom_ = this.element_.querySelector('.like-cell-bottom');
this.likeCellFooter_ = this.element_.querySelector('.like-cell-footer');
this.likeContent_ = this.element_.querySelector('.like-content');
this.likeContentWrapper_ = this.element_.querySelector('.like-content-wrapper');
this.likeHeader_ = this.element_.querySelector('.like-header');
this.likeLoadMore_ = this.element_.querySelector('.like-load-more');
this.likePrice_ = this.element_.querySelector('.like-price');
this.likeQuantities_ = this.element_.querySelector('.like-quantities');
this.likeRow_ = this.element_.querySelector('.like-row');
this.likeTitle_ = this.element_.querySelector('.like-title');
this.likeTitleBlack_ = this.element_.querySelector('.like-title-black');
this.likeTitleRed_ = this.element_.querySelector('.like-title-red');
this.likeWrapper_ = this.element_.querySelector('.like-wrapper');
this.modalHotRefresh_ = this.element_.querySelector('.modal_hot-refresh');
this.modalHotRefreshIcon_ = this.element_.querySelector('.modal_hot-refresh-icon');
this.modalHotTags_ = this.element_.querySelector('.modal_hot-tags');
this.modalTitle_ = this.element_.querySelector('.modal_title');
this.nodeInserted_ = this.element_.querySelector('.node-inserted');
this.refresh_ = this.element_.querySelector('.refresh');
this.search_ = this.element_.querySelector('.search');
this.searchModalAutocomplete_ = this.element_.querySelector('.search-modal_autocomplete');
this.searchModalClear_ = this.element_.querySelector('.search-modal_clear');
this.searchModalClearIcon_ = this.element_.querySelector('.search-modal_clear-icon');
this.searchModalHistory_ = this.element_.querySelector('.search-modal_history');
this.searchModalHistorytags_ = this.element_.querySelector('.search-modal_historytags');
this.searchModalHot_ = this.element_.querySelector('.search-modal_hot');
this.searchModalTitle_ = this.element_.querySelector('.search-modal_title');
this.topClear_ = this.element_.querySelector('.top_clear');
this.topClearIcon_ = this.element_.querySelector('.top_clear-icon');
this.topInput_ = this.element_.querySelector('.top_input');
this.topInputText_ = this.element_.querySelector('.top_input-text');
this.topSearch_ = this.element_.querySelector('.top_search');
this.topSearchIcon_ = this.element_.querySelector('.top_search-icon');
 */