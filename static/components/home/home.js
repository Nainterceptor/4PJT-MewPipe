(function () {
    "use strict";
    angular.module('mewpipe.home', [
        'mewpipeServices'
    ])
        .controller('HomeController', ['mediaFactory', 'paginationFactory', '$location', HomeController])
        .filter('emptyToEnd', [EmptyToEnd])
    ;

    function EmptyToEnd() {
        return function (array, key) {
            if (!angular.isArray(array)) return;
            var present = array.filter(function (item) {
                return item[key];
            });
            var empty = array.filter(function (item) {
                return !item[key]
            });
            return present.concat(empty);
        };
    }

    function HomeController(mediaFactory, paginationFactory, $location) {
        var me = this;
        this.mostShared = [];

        mediaFactory.getMediasByViews().success(function (response) {
            me.baseUrl = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/";
            me.medias = response;
            paginationFactory.setPagination(me.medias);
            me.page = paginationFactory.getParams();
        });

        mediaFactory.getMediasByShares().success(function (response) {
            me.mediasShares = response;
            paginationFactory.setPagination(me.mediasShares);
            me.pageShare = paginationFactory.getParams();
        });
    }
}());