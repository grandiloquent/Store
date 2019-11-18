;(function () {
    var Categories = function Categories() {

    };
    Categories.prototype.initialize = function () {
        this.element_ = document.getElementById('category-container');
        if (!this.element_) return;
        this.tabItemContainers_ = this.element_.querySelectorAll('.tab-item-container');
        this.tabContainer_ = this.element_.querySelector('.tab-container');
        this.calculateSize();
        this.setupTabItems();

    };
    Categories.prototype.calculateSize = function () {
        this.height_ = window.innerHeight;
        this.scorllHeight_ = this.tabContainer_.scrollHeight;
    };
    Categories.prototype.setupTabItems = function () {
        var that = this;
        this.tabItemContainers_.forEach(function (element) {
            element.addEventListener('click', function (event) {
                var cur = event.currentTarget;
                var selected = that.element_.querySelector('.tab-item-container.selected');
                if (cur !== selected) {
                    selected.classList.remove('selected');
                    cur.classList.add('selected');
                    var bottom = cur.getClientRects()[0].bottom;
                    if (bottom > window.innerHeight / 2) {
                        var offset = that.tabContainer_.scrollHeight - that.tabContainer_.offsetHeight - that.tabContainer_.scrollTop
                        var scroll = Math.min(bottom - window.innerHeight / 2, offset);
                        console.log(offset, scroll,bottom);
                        that.tabContainer_.scrollBy(0, scroll);

                    }
                }
            })
        });
    };
    var categories = new Categories();
    categories.initialize();
})();