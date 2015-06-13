(function () {
    "use strict";
    angular.module('mewpipe.adminPanel', [])
        .controller('AdminPanelController', ['userFactory', 'mediaFactory','notificationFactory', '$location', AdminPanelController])
        .directive('users', ['userFactory', 'paginationFactory', UsersDirective])
        .directive('medias', ['mediaFactory', 'paginationFactory', MediasDirective])
        .directive('modalUpdateUser', ['userFactory', ModalUpdateUserDirective])
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

    function AdminPanelController(userFactory, mediaFactory, notificationFactory, $location) {
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        var me = this;
        this.media = mediaFactory;
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

        userFactory.getUsers().success(function(response){
            me.users = response;
        });

        mediaFactory.getMedias().success(function (response) {
            me.medias = response;
        });
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
                this.user = userFactory;
                this.updateUser = function (user) {
                    angular.element('#updateUserModal' + user.id).appendTo('body').modal('show');
                };
                userFactory.getUsers().success(function(response){
                    me.users = response;
                    paginationFactory.setPagination(me.users);
                    me.page = paginationFactory.getParams();
                });
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
                this.media = mediaFactory;
                mediaFactory.getMedias().success(function (response) {
                    me.medias = response;
                    paginationFactory.setPagination(me.medias);
                    me.page = paginationFactory.getParams();
                });
            }
        }
    }

    function ModalUpdateUserDirective(userFactory) {
        return {
            restrict: 'E',
            templateUrl: 'app/templates/update-user.html',
            controllerAs: 'updateUser',
            controller: function ($scope, $element, $attrs) {
                this.update = function () {
                    userFactory.updateUser(this.id, this.email, this.firstname, this.lastname, this.nickname, this.password);
                    angular.element('#updateUserModal').appendTo('body').modal('hide');
                };
            }
        }
    }
}());