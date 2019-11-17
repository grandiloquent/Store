;(function () {
    var SearchResults = function SearchResults() {
    };
    SearchResults.prototype.initialize = function () {
        this.element_ = document.getElementById('results-header');
        if (!this.element_) return;
        this.filterGroupItems_ = Array.prototype.slice.call(
            document.querySelectorAll('.filter_group-item'));
        this.filter_ = this.filterGroupItems_.pop();
        this.searchBar_ = this.element_.querySelector('.search-bar');
        this.filterGroup_ = this.element_.querySelector('.filter_group');
        this.filterTags_ = this.element_.querySelector('.filter_tags');
        this.searchList_ = document.querySelector('.search-list');

        this.setupFilters();
        this.setupFilter();
        this.calculateSize();
    };
    SearchResults.prototype.setupFilters = function () {
        var elements = this.filterGroupItems_;
        for (var i = 0; i < elements.length; i++) {
            elements[i].addEventListener('click', function (event) {
                var cur = event.currentTarget;
                for (var j = 0; j < elements.length; j++) {
                    if (elements[j] !== cur) {
                        elements[j].classList.remove('active');
                    } else {
                        elements[j].classList.add('active');
                    }
                }
            })
        }
    };
    SearchResults.prototype.setupFilter = function () {
        if (!this.filterMask_) {
            this.filterMask_ = document.querySelector('.filter-mask');
        }
        if (!this.filterMain_) {
            this.filterMain_ = document.querySelector('.filter-main');
        }
        if (!this.submit_) {
            this.submit_ = document.querySelector('.submit');
        }
        this.filter_.addEventListener('click', function () {
            this.filterMask_.classList.remove('hidden');
            this.filterMain_.classList.remove('hidden');
        }.bind(this));
        this.submit_.addEventListener('click', function () {
        }.bind(this));
        this.filterMask_.addEventListener('click', this.onCloseFilter.bind(this));
        this.submit_.addEventListener('click', this.onCloseFilter.bind(this));
        /*
this.empty_ = this.element_.querySelector('.empty');
this.filterBody_ = this.element_.querySelector('.filter-body');
this.filterFooter_ = this.element_.querySelector('.filter-footer');
this.filterFooterButton_ = this.element_.querySelector('.filter-footer-button');
this.filterFooterButtonGroup_ = this.element_.querySelector('.filter-footer-button-group');
this.filterHeader_ = this.element_.querySelector('.filter-header');
this.filterHeaderTitle_ = this.element_.querySelector('.filter-header-title');
this.filterMain_ = this.element_.querySelector('.filter-main');
this.hidden_ = this.element_.querySelector('.hidden');
this.submit_ = this.element_.querySelector('.submit');
this.title_ = this.element_.querySelector('.title');
if(!this.empty_ ){this.empty_ =document.querySelector('.empty');}
if(!this.filterBody_ ){this.filterBody_ =document.querySelector('.filter-body');}
if(!this.filterFooter_ ){this.filterFooter_ =document.querySelector('.filter-footer');}
if(!this.filterFooterButton_ ){this.filterFooterButton_ =document.querySelector('.filter-footer-button');}
if(!this.filterFooterButtonGroup_ ){this.filterFooterButtonGroup_ =document.querySelector('.filter-footer-button-group');}
if(!this.filterHeader_ ){this.filterHeader_ =document.querySelector('.filter-header');}
if(!this.filterHeaderTitle_ ){this.filterHeaderTitle_ =document.querySelector('.filter-header-title');}
if(!this.hidden_ ){this.hidden_ =document.querySelector('.hidden');}
if(!this.title_ ){this.title_ =document.querySelector('.title');}
*/
    };
    SearchResults.prototype.calculateSize = function () {

        var transformStyle = 'transform: translate3d(0px, -' + (
            this.searchBar_.clientHeight +
            this.filterGroup_.clientHeight
        ) + 'px, 0px);';
        var paddingStyle = 'padding-top:' + this.filterTags_.clientHeight + 'px;';
        var offset = this.searchList_.querySelector('.skelecton-item').clientHeight;
        var that = this;
        window.addEventListener('scroll', function () {
            if (window.pageYOffset > offset) {
                that.element_.classList.add('fixed');
                that.element_.style = transformStyle;
                that.searchList_.style = paddingStyle;
            } else {
                that.element_.classList.remove('fixed');
                that.element_.removeAttribute('style');
                that.searchList_.removeAttribute('style');
            }
        })
    };
    SearchResults.prototype.onCloseFilter = function () {
        this.filterMain_.classList.add('hidden');
        this.filterMask_.classList.add('hidden');
    };
    var searchReuslts = new SearchResults();
    searchReuslts.initialize();
})();