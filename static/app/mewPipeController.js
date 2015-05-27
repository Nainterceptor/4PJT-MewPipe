(function () {
    "use strict";
    angular.module('mewpipe')
        .controller('MainController', ['$router', '$scope', 'userFactory', 'notificationFactory', 'themesFactory', MainController])
        .controller('AuthenticationController', ['userFactory', 'notificationFactory', AuthenticationController])
        .directive('modalSignIn', ['userFactory', 'notificationFactory', ModalSignInDirective])
        .directive('modalSignUp', ['userFactory', 'notificationFactory', ModalSignUpDirective])
        .directive('modalUpdateUser', ['userFactory', ModalUpdateUserDirective])
        .directive('pagination', ['paginationFactory', PaginationDirective])
    ;

    function MainController($router, $scope, userFactory, notificationFactory, themesFactory) {
        var me = this;
        this.themes = themesFactory.themes;
        $router.config([
            {path: '/', component: 'home'},
            {path: '/player/:id', component: 'player'},
            {path: '/dashboard', component: 'dashboard'},
            {path: '/dashboard/manage-video', component: 'dashboard'},
            {path: '/dashboard/profile', component: 'dashboard'},
            {path: '/admin-panel', component: 'adminPanel'}
        ]);
        this.chooseTheme = function (theme) {
            me.theme = themesFactory.themes[theme];
            themesFactory.saveTheme(theme);
        };
        me.theme = themesFactory.getTheme() ? themesFactory.themes[themesFactory.getTheme()] : themesFactory.themes['default'];
        this.alerts = notificationFactory.alerts;
        $scope.$on('alert:updated', function () {
            $scope.$apply();
        });
        this.closeAlert = function (index) {
            notificationFactory.delAlert(index);
        };
        this.user = userFactory;
    }

    function AuthenticationController(userFactory, notificationFactory) {
        this.logIn = function () {
            angular.element('#signInModal').appendTo('body').modal('show');
        };
        this.logOut = function () {
            userFactory.logOut();
        };
        this.signUp = function () {
            angular.element('#signUpModal').appendTo('body').modal('show');
        };
        this.getUsers = function () {
            userFactory.getUsers();
        };
        this.getUser = function () {
            userFactory.getUser();
        };
        this.updateUser = function () {
            angular.element('#updateUserModal').appendTo('body').modal('show');
        };
        this.deleteUser = function() {
            userFactory.deleteUser(this.id);
        };
    }

    function ModalSignInDirective(userFactory) {
        return {
            restrict: 'E',
            templateUrl: 'app/templates/sign-in-modal.html',
            controllerAs: 'signIn',
            controller: function () {
                this.logIn = function () {
                    userFactory.logIn(this.email, this.password);
                    angular.element('#signInModal').appendTo('body').modal('hide');
                };
            }
        }
    }

    function ModalSignUpDirective(userFactory, notificationFactory) {
        return {
            restrict: 'E',
            templateUrl: 'app/templates/sign-up-modal.html',
            controllerAs: 'signUp',
            controller: function () {
                this.signUp = function () {
                    userFactory.signUp(this.email, this.nickname, this.password);
                    angular.element('#signUpModal').appendTo('body').modal('hide');
                };
            }
        }
    }

    function ModalUpdateUserDirective(userFactory) {
        return {
            restrict: 'E',
            templateUrl: 'app/templates/update-user.html',
            controllerAs: 'updateUser',
            controller: function () {
                this.update = function () {
                    userFactory.updateUser(this.id, this.email, this.firstname, this.lastname, this.nickname, this.password);
                    angular.element('#updateUserModal').appendTo('body').modal('hide');
                };
            }
        }
    }

    function PaginationDirective(paginationFactory) {
        return {
            restrict: 'E',
            templateUrl: 'app/templates/pagination.html',
            controllerAs: 'pagination',
            controller: function () {
                this.page = {
                    currentPage: paginationFactory.getParams().currentPage,
                    numPerPage: paginationFactory.getParams().numPerPage,
                    numberOfPages: paginationFactory.numberOfPages()
                };
            }
        }
    }
}());