;(function () {

    var Details = function Details() {

    };
    Details.prototype.hideAttributeMask = function () {
        this.detailAttributeMask_.style.display = 'block';
        this.detailAttributeMask_.style.zIndex = '100';
        this.detailAttributeMask_.style.opacity = '0';

        var speed = 250;
        var start = +new Date;
        var from = 1;
        var to = 0;

        var element = this.detailAttributeMask_;

        var timer = setInterval(function () {

            var timeElap = +new Date - start;

            if (timeElap > speed) {
                element.style.display = 'none';
                element.style.opacity = '0';

                clearInterval(timer);
                return;

            }
            console.log((((to - from) * (Math.floor((timeElap / speed) * 100) / 100)) + from) + '');

            element.style.opacity = (((to - from) * (Math.floor((timeElap / speed) * 100) / 100)) + from) + '';

        }, 4);

    };
    Details.prototype.showAttributeMask = function () {
        this.detailAttributeMask_.style.display = 'block';
        this.detailAttributeMask_.style.zIndex = '100';
        this.detailAttributeMask_.style.opacity = '0';

        var speed = 250;
        var start = +new Date;
        var from = 0;
        var to = 1;

        var element = this.detailAttributeMask_;

        var timer = setInterval(function () {

            var timeElap = +new Date - start;

            if (timeElap > speed) {

                element.style.opacity = '1';

                clearInterval(timer);
                return;

            }
            console.log((((to - from) * (Math.floor((timeElap / speed) * 100) / 100)) + from) + '');

            element.style.opacity = (((to - from) * (Math.floor((timeElap / speed) * 100) / 100)) + from) + '';

        }, 4);

    };
    Details.prototype.initialize = function () {
        this.element_ = document.getElementById('details-content');
        if (!this.element_) return;
        this.detailAttribute_ = this.element_.querySelector('.detail-attribute');
        this.detailAttributeMask_ = document.querySelector('.detail-attribute-mask');

        /*
        this.detailAttributeMask_ = this.element_.querySelector('.detail-attribute-mask');

        if(!this.detailAttributeMask_ ){this.detailAttributeMask_ =document.querySelector('.detail-attribute-mask');}

        */

        this.adjustSize();
        this.detailAttribute_.addEventListener('click', function () {
            this.showAttributeMask();
        }.bind(this));
        this.detailAttributeMask_.addEventListener('click', function () {
            this.hideAttributeMask();
        }.bind(this));
    };
    Details.prototype.adjustSize = function () {
        document.documentElement.style.fontSize = "100%";
    };
    var details = new Details();
    details.initialize();
})();