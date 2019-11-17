;(function () {

    var Details = function Details() {

    };
    Details.prototype.initialize = function () {
        this.adjustSize();
    };
    Details.prototype.adjustSize = function () {
        document.documentElement.style.fontSize = "100%";
    };
    var details = new Details();
    details.initialize();
})();