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
    };
// 计算相应的尺寸
    Home.prototype.calculateSize = function () {
        var gutter = (window.innerWidth * 0.01333).toFixed(2);
        var cellWidth = window.innerWidth / 2 - gutter / 2;
        this.cellWidth_ = 'width:' + cellWidth + 'px';
        this.cellSize_ = 'width:' + cellWidth + 'px;height:' + (cellWidth + 96) + 'px';
        this.imgSize_ = 'width:' + cellWidth + 'px;height:' + cellWidth + 'px';
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
        this.setupItems();
    };
    Home.prototype.onRefreshFailed = function (error) {
        console.log(error);
    };
    Home.prototype.onRefreshSuccess = function (obj) {
        that.offset_ *= 2;
        var content = [];
        for (var i = 0; i < obj.length; i++) {
            var pattern = "<div class=\"hot-tag\">" + obj[i] + "</div>";
            content.push(pattern);
        }
        this.modalHotTags_.innerHTML = content.join('');
        this.setupKeyword();
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
    Home.prototype.setupItems = function () {
        document.querySelectorAll('.like-cell')
            .forEach(function (element) {
                element
                    .addEventListener('click', function (event) {
                        window.location = "/store/details/" + event.currentTarget.getAttribute('data-id')
                    });
            });
    };
    // Home.prototype.setupKeyword = function () {
    //     var that = this;
    //     this.modalHotTags_
    //         .querySelectorAll('.hot-tag')
    //         .forEach(function (element) {
    //             element.addEventListener('click', function (event) {
    //                 window.location = "/store/results?keyword=" + event.currentTarget.textContent;
    //                 that.addHistory(event.currentTarget.textContent);
    //             })
    //         });
    // };

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
    // 显示搜索框
    Home.prototype.showSearchModal = function () {
        if (!this.searchModal_) {
            this.searchModal_ = new SearchModal();
            this.searchModal_.initialize(function () {
                this.swiperPagination_.style.display = 'flex';
            }.bind(this));
        }
        this.searchModal_.show();
        this.swiperPagination_.style.display = 'none';
    };
    var home = new Home();
    home.initialize();
})();