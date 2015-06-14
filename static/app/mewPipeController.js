(function () {
    "use strict";
    angular.module('mewpipe')
        .controller('MainController', ['$router', '$scope', 'userFactory', 'notificationFactory', 'themesFactory','$cookies', MainController])
        .controller('AuthenticationController', ['userFactory', 'notificationFactory', AuthenticationController])
        .directive('modalSignIn', ['userFactory', 'notificationFactory', ModalSignInDirective])
        .directive('modalSignUp', ['userFactory', 'notificationFactory', ModalSignUpDirective])
    ;

    function MainController($router, $scope, userFactory, notificationFactory, themesFactory, $cookies) {
        var me = this;
        if ($cookies.get('userId')) {
            userFactory.user = {
                id: $cookies.get('userId')
            };
            userFactory.getUser()
                .success(function (response) {
                    userFactory.setUser(response);
                })
                .error(function (response) {
                    console.log(response);
                });
        }
        this.themes = themesFactory.themes;
        $router.config([
            {path: '/', component: 'home'},
            {path: '/player/:id/', component: 'player'},
            {path: '/dashboard/', component: 'dashboard'},
            {path: '/account/', component: 'account'},
            {path: '/manage-video/', component: 'manageVideo'},
            {path: '/admin-panel/', component: 'adminPanel'},
            {path: '/admin-panel/medias/', component: 'adminPanel'},
            {path: '/admin-panel/users/', component: 'adminPanel'},
            {path: '/login/', component: 'login'},
            {path: '/login/twitter/', component: 'login'},
            {path: '/login/twitter/callback/', component: 'login'},
            {path: '/upload/', component: 'upload'},
            {path: '/update-video/:id', component: 'updateVideo'},
            {path: '/update-user/:id', component: 'upload'}
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
        this.logOut = function () {
            userFactory.logOut();
        };
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
        this.deleteUser = function () {
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
}());