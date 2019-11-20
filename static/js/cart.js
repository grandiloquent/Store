;(function () {
    var Cart = function Cart() {

    };
    // ==============================================
    Cart.prototype.initialize = function () {
        this.element_ = document.getElementById('cart');
        if (!this.element_) return;
        this.adjustSize();
        this.loadCartFromStorage();
    };
    Cart.prototype.loadCartFromStorage = function () {
        if (!window.localStorage) return;

        var cart = localStorage.getItem("cart")
        if (!cart) return;
        this.items_ = JSON.parse(cart);
    };
    Cart.prototype.adjustSize = function () {
        document.documentElement.style.fontSize = "100%";
    };

    // ==============================================
    var cart = new Cart();
    cart.initialize();

})();
