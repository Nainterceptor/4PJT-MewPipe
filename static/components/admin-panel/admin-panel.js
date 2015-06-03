(function () {
    "use strict";
    angular.module('mewpipe.adminPanel', [])
        .controller('AdminPanelController', ['userFactory', 'notificationFactory', '$location', AdminPanelController])
        .directive('users', ['userFactory', 'paginationFactory', UsersDirective])
        .directive('medias', ['mediaFactory', 'paginationFactory', MediasDirective])
        .filter('startFrom', AdminPanelFilter)
    ;

    function AdminPanelFilter() {
        return function (input, start) {
            if (!input || !input.length) {
                return;
            }
            start = +start; //parse to int
            return input.slice(start);
        }
    }

    function AdminPanelController(userFactory, notificationFactory, $location) {
        this.canActivate = function(){
            if (!userFactory.accessToken){
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        var me = this;
        this.user = userFactory;

        if ($location.url() != '/admin-panel/medias') {
            this.activeTab = 'users';
        } else {
            this.activeTab = 'medias';
        }

        this.active = function (tab) {
            if (tab === 'users') {
                $location.url('/admin-panel/users');
            } else {
                $location.url('/admin-panel/medias');
            }
            me.activeTab = tab;
        };
    }

    function UsersDirective(userFactory, paginationFactory) {
        return {
            restrict: 'E',
            templateUrl: 'components/admin-panel/users.html',
            scope: true,
            bindToController: true,
            controllerAs: 'users',
            controller: function ($scope, $element, $attrs) {
                var me = this;
                this.user = userFactory.users;
                paginationFactory.setPagination(me.user, 0, 5);
            }
        }
    }

    function MediasDirective(mediaFactory, paginationFactory) {
        return {
            restrict: 'E',
            templateUrl: 'components/admin-panel/medias.html',
            scope: true,
            bindToController: true,
            controllerAs: 'medias',
            controller: function ($scope, $element, $attrs) {
                var me = this;
                this.media = mediaFactory.medias;
                paginationFactory.setPagination(me.media, 0, 5);
            }
        }
    }
}());