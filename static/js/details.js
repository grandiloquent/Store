;(function () {

    var Details = function Details() {

    };

    // ==============================================

    Details.prototype.hideAttributeMask = function () {

        function animate(time) {
            requestAnimationFrame(animate);
            TWEEN.update(time);
        }

        requestAnimationFrame(animate);

        var element = this.detailAttributeMask_;
        var start = {s: 1};
        var tween = new TWEEN.Tween(start) // Create a new tween that modifies 'coords'.
            .to({s: 0}, 400) // Move to (300, 200) in 1 second.
            .easing(TWEEN.Easing.Quadratic.Out) // Use an easing function to make the animation smooth.
            .onUpdate(() => {
                element.style.opacity = start.s + '';
            }).onComplete(function () {
                element.style.display = 'none';
            })

            .start(); // Start the tween immediately.

    };
    Details.prototype.showAttributeMask = function () {
        this.detailAttributeMask_.style.display = 'block';
        this.detailAttributeMask_.style.zIndex = '100';
        this.detailAttributeMask_.style.opacity = '0';

        function animate(time) {
            requestAnimationFrame(animate);
            TWEEN.update(time);
        }

        requestAnimationFrame(animate);

        var element = this.detailAttributeMask_;
        var start = {s: 0};
        var tween = new TWEEN.Tween(start) // Create a new tween that modifies 'coords'.
            .to({s: 1}, 400) // Move to (300, 200) in 1 second.
            .easing(TWEEN.Easing.Quadratic.Out) // Use an easing function to make the animation smooth.
            .onUpdate(() => {
                element.style.opacity = start.s + '';
            })

            .start(); // Start the tween immediately.

    };
    Details.prototype.initialize = function () {
        this.element_ = document.getElementById('details-content');
        if (!this.element_) return;
        this.detailAttribute_ = this.element_.querySelector('.detail-attribute');
        this.detailAttributeMask_ = document.querySelector('.detail-attribute-mask');
        this.detailAttributeContent_ = this.detailAttributeMask_.querySelector('.detail-attribute-content');
        this.detailAttributeHeader_ = this.detailAttributeMask_.querySelector('.detail-attribute-header');

        /*
        this.detailAttributeMask_ = this.element_.querySelector('.detail-attribute-mask');

        if(!this.detailAttributeMask_ ){this.detailAttributeMask_ =document.querySelector('.detail-attribute-mask');}

        */
        this.setupSwipe();
        this.setupGoBack();
        this.adjustSize();
        this.setupScroll();
        // ==============================================
        this.detailAttribute_.addEventListener('click', function () {
            this.showAttributeMask();
        }.bind(this));
        this.detailAttributeMask_.addEventListener('click', function () {
            this.hideAttributeMask();
        }.bind(this));
        this.detailAttributeHeader_.addEventListener('click', function (event) {
            event.stopPropagation();
        });
        this.detailAttributeContent_.addEventListener('click', function (event) {
            event.stopPropagation();
            event.stopImmediatePropagation();
        });

    };
    Details.prototype.adjustSize = function () {
        document.documentElement.style.fontSize = "100%";
    };
    Details.prototype.setupSwipe = function () {
        this.swipe_ = this.element_.querySelector('.swipe');
        this.imageSlideCurrentIndex_ = this.element_.querySelector('.image-slide-current-index');
        this.imageSlideCount_ = this.element_.querySelector('.image-slide-count');

        this.imageSlideCount_.textContent = this.element_.querySelectorAll('.swipe-wrap>div').length + '';
        var that = this;
        this.swipe_.style.height = this.swipe_.style.width = window.innerWidth + 'px';
        new Swipe(this.swipe_, {
            auto: 3000, callback: function (index, elem) {
                that.imageSlideCurrentIndex_.textContent = (index + 1) + '';
            },
        });
    };
    Details.prototype.setupGoBack = function () {
        this.btnGoBack_ = this.element_.querySelector('.btn-go-back');
        this.btnGoBack_.addEventListener('click', function () {
            window.location = "/store";
        });
    };
    Details.prototype.setupScroll = function () {
        this.detailHeaderLayout_ = this.element_.querySelector('.detail-header-layout');
        var height = this.detailHeaderLayout_.getClientRects()[0].height;
        this.detailHeaderLayout_.style.display = 'none';
        var that = this;
        window.addEventListener('scroll', function () {
            if (window.pageYOffset > height) {
                that.detailHeaderLayout_.style.display = 'block';
            } else {
                that.detailHeaderLayout_.style.display = 'none';
            }
        });
    };
    // ==============================================

    var details = new Details();
    details.initialize();
})();