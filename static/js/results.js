;(function () {

    var Results = function Results() {


    };

    Results.prototype.initialize = function () {
        this.listColumn_ = document.querySelector('.list_column');

        this.loadImages();

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

    var results = new Results();
    results.initialize();

})();