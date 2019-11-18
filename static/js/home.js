;(function () {
    'use strict';
    var Home = function Home() {
    };
    Home.prototype.adjustSize = function () {
        var that = this;
        _.every(this.element_, '.like-cell', function (element) {
            element.style = that.cellSize_;
        });
        _.every(this.element_, '.like-cell img', function (element) {
            element.style = that.imgSize_;
        });
        _.every(this.element_, '.like-cell-footer', function (element) {
            element.style = that.cellWidth_;
        });
    }
// 计算相应的尺寸
    Home.prototype.calculateSize = function () {
        var gutter = (window.innerWidth * 0.01333).toFixed(2);
        var cellWidth = window.innerWidth / 2 - gutter / 2;
        this.cellWidth_ = 'width:' + cellWidth + 'px';
        this.cellSize_ = 'width:' + cellWidth + 'px;height:' + (cellWidth + 96) + 'px';
        this.imgSize_ = 'width:' + cellWidth + 'px;height:' + cellWidth + 'px';
    };
    Home.prototype.fetchHistory = function () {
        //"<div class=\"history-tag\">"+history+"</div>"
        if (!window.localStorage) return;
        var searchHistory = window.localStorage.getItem('search-history');
        if (!searchHistory) return;
        var obj = JSON.parse(searchHistory);
        var content = [];
        for (var i = 0; i < obj.length; i++) {
            var pattern =
                "<div class=\"history-tag\">" + obj[i] + "</div>"
            content.push(pattern);
        }
        if (!content.length) return;
        if (!this.searchModalHistorytags_) {
            this.searchModalHistorytags_ = document.querySelector('.search-modal_historytags');
        }
        this.searchModalHistorytags_.innerHTML = content.join('');
    };
    Home.prototype.initialize = function () {
        this.element_ = document.getElementById('like-content');
        if (!this.element_) {
            console.log('cant detect the base element. stop rendering the home page.');
            return;
        }
        this.calculateSize();
        this.adjustSize();
        this.setupSlide();
        this.setupHomeSearch();
    };
    Home.prototype.onRefreshFailed = function (error) {
        console.log(error);
    };
    Home.prototype.setupSlide = function () {
        this.swiperActiveIndex_ = document.getElementById('swiper-active-index');
        this.swiperPagination_ = document.querySelector('.swiper-pagination');


        var that = this;
        new Swipe(document.getElementById('slider'), {
            startSlide: 2,
            speed: 400,
            auto: 3000,
            continuous: true,
            disableScroll: false,
            stopPropagation: false,
            callback: function (index, elem) {
                that.swiperActiveIndex_.innerHTML = index + 1;
            },
            transitionEnd: function (index, elem) {
            }
        });
    };
    Home.prototype.onRefreshSuccess = function (obj) {
        var content = [];
        for (var i = 0; i < obj.length; i++) {
            var pattern = "<div class=\"hot-tag\">" + obj[i] + "</div>";
            content.push(pattern);
        }
        this.modalHotTags_.innerHTML = content.join('');
        this.setupKeyword();
    };
    Home.prototype.addHistory = function (keywords) {
        if (!window.localStorage) return;
        var obj = (JSON.parse(window.localStorage.getItem('search-history'))) || [];
        if (obj.indexOf(keywords) === -1) {
            obj.push(keywords);
        }
        window.localStorage.setItem('search-history', JSON.stringify(obj));
    };
    Home.prototype.setupHomeSearch = function () {
        var elements = document.querySelectorAll('.home-search');
        var top = elements[0].getClientRects()[0].top + window.pageYOffset;
        window.addEventListener('scroll', function () {
            if (window.pageYOffset > top) {
                elements[1].classList.remove('hide');
            } else {
                elements[1].classList.add('hide');
            }
        });
        for (var i = 0; i < elements.length; i++) {
            _.click(elements[i], this.showSearchModal.bind(this));
        }
    };
    Home.prototype.setupKeyword = function () {
        var that = this;
        this.modalHotTags_
            .querySelectorAll('.hot-tag')
            .forEach(function (element) {
                element.addEventListener('click', function (event) {
                    window.location = "/store/results?keyword=" + event.currentTarget.textContent;
                    that.addHistory(event.currentTarget.textContent);
                })
            });
    };
    Home.prototype.setupRefresh = function () {
        if (!this.modalHotRefresh_)
            this.modalHotRefresh_ = document.querySelector('.modal_hot-refresh');
        if (!this.modalHotTags_)
            this.modalHotTags_ = document.querySelector('.modal_hot-tags');
        var that = this;
        this.modalHotRefresh_.addEventListener('click', function () {
            fetch("/store/api/search?method=fetch")
                .then(function (response) {
                    return response.json();
                })
                .then(that.onRefreshSuccess.bind(that))
                .catch(that.onRefreshFailed)
        });
    };
    // 显示搜索框
    Home.prototype.showSearchModal = function () {
        if (!this.searchModal_) {
            this.searchModal_ = new SearchModal();
            this.searchModal_.initialize(function () {
                this.swiperPagination_.style.display = 'flex';
            }.bind(this));
            this.setupRefresh();
            this.setupKeyword();
            this.searchModalHistorytags_ = document.querySelector('.search-modal_historytags');


            if (!this.searchModalClear_) {
                this.searchModalClear_ = document.querySelector('.search-modal_clear');
                var that = this;
                this.searchModalClear_.addEventListener('click', function () {
                    window.localStorage.removeItem('search-history');
                    that.searchModalHistorytags_.innerHTML = '';
                });
            }

        }
        this.searchModal_.show();
        this.fetchHistory();
        this.swiperPagination_.style.display = 'none';
    };
    var home = new Home();
    home.initialize();
})();