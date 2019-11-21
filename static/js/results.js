;(function () {

    var Results = function Results() {

    };

    Results.prototype.initialize = function () {

        this.element_ = document.getElementById('results-header');
        if (!this.element_) return;
        this.filterGroupItems_ = Array.prototype.slice.call(
            document.querySelectorAll('.filter_group-item'));
        this.filter_ = this.filterGroupItems_.pop();
        this.searchBar_ = this.element_.querySelector('.search-bar');
        this.filterGroup_ = this.element_.querySelector('.filter_group');
        this.filterTags_ = this.element_.querySelector('.filter_tags');
        this.searchList_ = document.querySelector('.search-list');


        this.listColumn_ = document.querySelector('.list_column');

        this.calculateSize();
        this.setupFilter();
        this.loadImages();

        this.setupGoBack();

        var review = document.getElementById('review');
        var quantities = document.getElementById('quantities');
        var price = document.getElementById('price');

        review.addEventListener('click', function () {
            if (review.classList.contains('active')) return;
            this.touchServer(1);
            review.classList.add('active');
            quantities.classList.remove('active');
            price.classList.remove('active');

        }.bind(this));
        quantities.addEventListener('click', function () {
            if (quantities.classList.contains('active')) return;
            this.touchServer(2);
            quantities.classList.add('active');

            review.classList.remove('active');
            price.classList.remove('active');

        }.bind(this));
        price.addEventListener('click', function () {
            if (price.classList.contains('active')) return;
            this.touchServer(3);
            price.classList.add('active');
            review.classList.remove('active');
            quantities.classList.remove('active');
        }.bind(this));

    };
    /*
    <a class="item-link" href="/store/details/{{uid}}">
    <div class="item-image"><img class="image_src" data-src="/store/static/pictures/{{thumbnail}}"></div>
    <div class="item-info">
        <div class="item-info_title"><span>{{title}}</span></div>
        <div class="item-info_count">
            <div class="count_price"><i>￥</i>{{price}}</div>
            <div class="count_vol">成交 {{quantities}} 笔</div>
        </div>
    </div>
</a>
     */

    Results.prototype.touchServer = function (sorttype) {
        var that = this;
        _.touchServer({
            uri: '/store/api/search?method=like&limit=20&offset=0&sorttype=' + sorttype + '&keyword=迷',
            success: function (obj) {
                that.render(obj);
            }
        })
    };

    Results.prototype.render = function (obj) {
        var buf = [];

        for (var i = 0; i < obj.length; i++) {
            var uid = obj[i]['uid'];
            var thumbnail = obj[i]['thumbnail'];
            var price = obj[i]['price'];
            var title = obj[i]['title'];
            var quantities = obj[i]['quantities'];

            var pattern = "<div class=\"skelecton-item\"><a class=\"item-link\" href=\"/store/details/" + uid + "\"><div class=\"item-image\"><img class=\"image_src\" data-src=\"/store/static/pictures/" + thumbnail + "\"/></div><div class=\"item-info\"><div class=\"item-info_title\"><span>" + title + "</span></div><div class=\"item-info_count\"><div class=\"count_price\"><i>￥</i>" + price + "</div><div class=\"count_vol\">成交" + quantities + "笔</div></div></div></a></div>";
            buf.push(pattern);
        }
        this.listColumn_.innerHTML = buf.join('');
        this.loadImages();


    };
    Results.prototype.loadImages = function () {
        _.every('.image_src', function (element) {
            element.addEventListener('load', function (event) {
                event.currentTarget.style.zIndex = '1';
                event.currentTarget.style.opacity = '1';
            });
            element.src = element.getAttribute('data-src');
        });
    };
    Results.prototype.calculateSize = function () {

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
    Results.prototype.setupFilter = function () {
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

    };
    Results.prototype.onCloseFilter = function () {
        this.filterMain_.classList.add('hidden');
        this.filterMask_.classList.add('hidden');
    };
    Results.prototype.setupGoBack = function () {
        this.barBack_ = this.element_.querySelector('.bar_back');
        this.barBack_.addEventListener('click', function () {
            if (window.history.state) {
                window.history.back();
            } else {
                window.location = '/store';
            }
        }.bind(this));
    };

    var results = new Results();
    results.initialize();

})();