;(function () {
    'use strict';
    var Home = function Home() {
    };


    Home.prototype.adjustSize = function () {
        var that = this;
        _.every(this.element_, '.like-cell', function (element) {
            element.style = that.cellSize_;
        });
        // _.every(this.element_, '.like-cell img', function (element) {
        //     element.style = that.imgSize_;
        // });
        // _.every(this.element_, '.like-cell-footer', function (element) {
        //     element.style = that.cellWidth_;
        // });
    };
// 计算相应的尺寸
    Home.prototype.calculateSize = function () {
        var gutter = (window.innerWidth * 0.01333).toFixed(2);
        var cellWidth = window.innerWidth / 2 - gutter / 2;
        this.cellWidth_ = 'width:' + cellWidth + 'px';
        this.cellSize_ = 'width:' + cellWidth + 'px;height:' + (cellWidth + 96) + 'px';
        this.imgSize_ = 'width:' + cellWidth + 'px;height:' + cellWidth + 'px';
    };

    Home.prototype.loadingMore = function () {
        _.touchServer({
            uri: '/store/api/store?limit=10&offset=' + this.offset_,
            success: this.onSuccess.bind(this),
            failed: this.onFailed.bind(this)
        })
    };
    Home.prototype.onFailed = function () {
        this.likeLoadMore_.textContent = '没有更多了';
    };
    Home.prototype.onSuccess = function (obj) {
        if (!obj) {
            this.likeLoadMore_.textContent = '没有更多了';
            return
        }
        var buf = [];

        for (var i = 0; i < obj.length; i += 2) {
            buf.push('<div class="like-row">');
            var uid = obj[i][0];
            var title = obj[i][1];
            var price = obj[i][2];
            var thumbnail = obj[i][3];
            var quantities = obj[i][4] || 0;
            var pattern = "<div class=\"like-cell\" style=\"" + this.cellSize_ + "\" data-id=\"" + uid + "\"><img src=\"/store/static/pictures/" + thumbnail + "\"/><div class=\"like-cell-footer\"><span>" + title + "</span><div class=\"like-cell-tags\"></div><div class=\"like-cell-bottom\"><span class=\"like-price\">￥" + price + "</span> <span class=\"like-quantities\">" + quantities + "</span></div></div></div>";
            buf.push(pattern);
            if (i + 1 < obj.length) {

                uid = obj[i + 1][0];
                title = obj[i + 1][1];
                price = obj[i + 1][2];
                thumbnail = obj[i + 1][3];
                quantities = obj[i + 1][4] || 0;
                pattern = "<div class=\"like-cell\" style=\"" + this.cellSize_ + "\" data-id=\"" + uid + "\"><img src=\"/store/static/pictures/" + thumbnail + "\"/><div class=\"like-cell-footer\"><span>" + title + "</span><div class=\"like-cell-tags\"></div><div class=\"like-cell-bottom\"><span class=\"like-price\">￥" + price + "</span> <span class=\"like-quantities\">" + quantities + "</span></div></div></div>";
                buf.push(pattern);
            }
            buf.push('</div>')
        }
        document.getElementById('like-content')
            .insertAdjacentHTML('beforeend', buf.join(''));
        this.isLoading_ = false;
        this.offset_ += 10;
        this.setupItems();
        /*
                               <div class="like-row">
                                    <div class="like-cell" data-id="{{uid}}">
                                        <img src="{{thumbnail}}">
                                        <div class="like-cell-footer">
                                            <span>{{title}}</span>
                                            <div class="like-cell-bottom">
                                                <span class="like-price">￥{{price}}</span>
                                                <span class="like-quantities">{{quantities}}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
           */
    };
    Home.prototype.setupScroll = function () {
        var element = document.querySelector('.like-row:last-child');
        this.offsetTop_ = element.offsetTop;
        var that = this;
        window.addEventListener('scroll', function () {
            if (!that.isLoading_ && window.pageYOffset > that.offsetTop_) {
                that.isLoading_ = true;
                that.loadingMore();
            }
        });
    };

    Home.prototype.initialize = function () {
        this.element_ = document.getElementById('like-content');
        if (!this.element_) {
            console.log('cant detect the base element. stop rendering the home page.');
            return;
        }
        this.likeLoadMore_ = document.querySelector('.like-load-more');

        // -----------------------------------
        this.offset_ = 10;
        this.calculateSize();
        this.adjustSize();
        this.setupSlide();
        this.setupHomeSearch();
        this.setupItems();
        this.setupScroll();
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
                element.onclick = function () {
                    window.location = "/store/details/" + event.currentTarget.getAttribute('data-id');
                }
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