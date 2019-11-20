;(function () {
    var OrderModal = function OrderModal() {

    };
    window['OrderModal'] = OrderModal;

    // ==============================================
    OrderModal.prototype.setupSubmit = function () {
        this.selectorBtn_ = this.element_.querySelector('.selector-btn');
        this.selectorBtn_.addEventListener('click', function () {
            this.onSubmit();
        }.bind(this));
    };
    OrderModal.prototype.onSubmit = function () {
        var phoneNumber = this.phoneNum_.value;
        if (!/^[0-9]{11}$/.test(phoneNumber)) {

            this.phoneNum_.setAttribute('placeholder', '请输入正确的手机号码');
            this.phoneNum_.value = '';
            this.phoneNum_.focus();
            return false;
        } else {
            window.localStorage && window.localStorage.setItem('phoneNumber', phoneNumber);
            this.touchServer(phoneNumber);
        }
    };
    OrderModal.prototype.onFailed = function () {
        this.hideOrderMask();
        _.toast("发生未知错误。<br>无法提交订单。<br>请稍后再试。");

    };
    OrderModal.prototype.onSuccess = function () {
        this.hideOrderMask();
        _.toast('已成功提交订单。');
    }
    OrderModal.prototype.touchServer = function (phoneNumber) {
        var obj = {
            'phoneNumber': phoneNumber,
            'quantities': this.amount_,
            'uid': this.element_.getAttribute('data-id'),
        };
        var that = this;
        _.touchServer({
            uri: '/',
            body: JSON.stringify(obj),
            success: function () {

            },
            failed: function (error) {
                that.onFailed(error);
            }
        });
    };

    OrderModal.prototype.setupPhoneNumber = function () {
        this.phoneNum_ = this.element_.querySelector('.phone-num-container input');
        if (window.localStorage && localStorage.getItem('phoneNumber')) {
            this.phoneNum_.value = localStorage.getItem('phoneNumber');
        }
    };
    OrderModal.prototype.initialize = function () {
        this.element_ = document.getElementById('detail-order-mask');
        if (!this.element_) return;
        this.amountDownBtn_ = this.element_.querySelector('.amount-down-btn');
        this.amountInput_ = this.element_.querySelector('.amount-input');
        this.amountUpBtn_ = this.element_.querySelector('.amount-up-btn');
        this.totalAmount_ = this.element_.querySelector('.total-amount .color-highlight');
        this.totalPrice_ = this.element_.querySelector('.total-price .price-num');
        this.detailOrderContainer_ = this.element_.querySelector('.detail-order-container');


        this.element_.addEventListener('click', function () {
            this.hideOrderMask();
        }.bind(this));
        this.element_.style.display = 'block';
        this.orderMaskHeight_ = this.detailOrderContainer_.getClientRects()[0].height;
        this.detailOrderContainer_.style.bottom = (this.orderMaskHeight_ * -1) + 'px';
        _.interceptClick(this.detailOrderContainer_);

        this.price_ = parseFloat(this.element_.getAttribute('data-price'));
        this.amount_ = 1;
        this.amountInput_.value = this.amount_ + '';

        this.setupControlButton(this.amountDownBtn_, true);
        this.setupControlButton(this.amountUpBtn_, false);
        this.setupInput();
        this.setupPhoneNumber();
        this.setupSubmit();
    };
    OrderModal.prototype.setupControlButton = function (element, isDown) {
        element.addEventListener('click', function () {
            if (isDown) {
                if (this.amount_ > 1)
                    this.amount_--;
                else return;
            } else {
                this.amount_++;
            }
            this.calculate(true);
        }.bind(this));
    };

    OrderModal.prototype.calculate = function (changeInput) {
        var amountText = this.amount_ + '';
        if (changeInput)
            this.amountInput_.value = amountText;
        this.totalAmount_.textContent = amountText;
        this.totalPrice_.textContent = (this.price_ * this.amount_).toFixed(2) + '';
    };
    OrderModal.prototype.setupInput = function () {
        this.amountInput_.addEventListener('input', function (event) {
            var amount = parseInt(event.currentTarget.value);
            if (isNaN(amount)) {
                this.amount_ = 1;
            } else {
                this.amount_ = amount;
            }
            this.calculate();
        }.bind(this));
    };
    OrderModal.prototype.showOrderMask = function () {
        this.element_.removeAttribute('style');

        function animate(time) {
            requestAnimationFrame(animate);
            TWEEN.update(time);
        }

        requestAnimationFrame(animate);
        var element = this.detailOrderContainer_;
        var start = {s: this.orderMaskHeight_ * -1};
        var tween = new TWEEN.Tween(start) // Create a new tween that modifies 'coords'.
            .to({s: 0}, 150) // Move to (300, 200) in 1 second.
            .easing(TWEEN.Easing.Linear.None) // Use an easing function to make the animation smooth.
            .onUpdate(function () {
                element.style.bottom = start.s + 'px';
            })
            .start(); // Start the tween immediately.
    };
    OrderModal.prototype.hideOrderMask = function () {
        function animate(time) {
            requestAnimationFrame(animate);
            TWEEN.update(time);
        }

        requestAnimationFrame(animate);
        var that = this;
        var element = this.detailOrderContainer_;
        var start = {s: 0};
        var tween = new TWEEN.Tween(start) // Create a new tween that modifies 'coords'.
            .to({s: this.orderMaskHeight_ * -1}, 150) // Move to (300, 200) in 1 second.
            .easing(TWEEN.Easing.Linear.None) // Use an easing function to make the animation smooth.
            .onUpdate(function () {
                element.style.bottom = start.s + 'px';
            })
            .onComplete(function () {
                that.element_.style.display = 'none';
            })
            .start(); // Start the tween immediately.

    };
    // ==============================================

})();

;(function () {

    var Details = function Details() {

    };

    // ==============================================
    ;
    Details.prototype.adjustSize = function () {
        document.documentElement.style.fontSize = "100%";
    }
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
            .onUpdate(function () {
                element.style.opacity = start.s + '';
            }).onComplete(function () {
                element.style.display = 'none';
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

        this.setupSwipe();
        this.setupGoBack();
        this.setupGoCart();
        this.setupOrderButton();
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
    Details.prototype.setupGoBack = function () {
        this.btnGoBack_ = this.element_.querySelector('.btn-go-back');
        this.btnGoBack_.addEventListener('click', function () {
            window.location = "/store";
        });
    };
    Details.prototype.setupGoCart = function () {
        this.btnGoCart_ = this.element_.querySelector('.btn-go-cart');
        this.btnGoCart_.addEventListener('click', function () {
            window.location = "/store/cart";
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
            .onUpdate(function () {
                element.style.opacity = start.s + '';
            })
            .start(); // Start the tween immediately.
    };

    Details.prototype.setupOrderButton = function () {
        this.orderButton_ = document.getElementById('order');
        this.orderButton_.addEventListener('click', function () {
            if (!this.orderModal) {
                this.orderModal = new OrderModal();
                this.orderModal.initialize();
            }
            this.orderModal.showOrderMask();
        }.bind(this));
    };


    // ==============================================

    var details = new Details();
    details.initialize();
})();



