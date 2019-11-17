;(function () {
    var Slider = function Slider() {

    };
    window['Slider'] = Slider;
    Slider.prototype.initialize = function () {
        this.element_ = document.getElementById('slider');
        if (!this.element_) return;
        /*

if(!this.swiperWrapper_ ){this.swiperWrapper_ =document.querySelector('.swiper-wrapper');}

*/
        this.swiperWrapper_ = this.element_.querySelector('.swiper-wrapper');

        this.delay_ = 3000;
        this.slide(0, -360);
    };
    Slider.prototype.next = function () {
        var that = this;
        setTimeout(function () {
            that.slide(-360, -720);
        }, this.delay_)
    }
    Slider.prototype.slide = function (from, to) {

        var speed = 400;

        var that = this;
        var start = +new Date;
        var timer = setInterval(function () {

            var timeElap = +new Date - start;
            if (timeElap > speed) {
                that.swiperWrapper_.style = 'transform: translate3d(' +
                    to + 'px'
                    + ', 0px, 0px); transition-duration: 0ms;';
                that.next();
                clearInterval(timer);
                return;
            }
            that.swiperWrapper_.style = 'transform: translate3d(' +
                (((to - from) * (Math.floor((timeElap / speed) * 100) / 100)) + from) + 'px'
                + ', 0px, 0px); transition-duration: 0ms;';
        });
    };
    var slider = new Slider();
    slider.initialize();
})();
