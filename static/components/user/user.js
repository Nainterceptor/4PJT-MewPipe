(function () {
    "use strict";
    angular.module('mewpipe.user', [])
        .controller('UserController', ['userFactory', 'mediaFactory', 'notificationFactory', '$routeParams', 'paginationFactory', '$location', UserController])
    ;

    function UserController(userFactory, mediaFactory, notificationFactory, $routeParams, paginationFactory, $location) {
        var me = this;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        userFactory.getOneUser($routeParams.id).success(function (response) {
            me.user = response;
        });
        mediaFactory.getOneUserMedias($routeParams.id).success(function (response) {
            me.baseUrl = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/";
            me.media = response;
            paginationFactory.setPagination(me.media);
            me.page = paginationFactory.getParams();
        });
    }
}());