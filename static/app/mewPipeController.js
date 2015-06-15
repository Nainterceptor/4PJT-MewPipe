(function () {
    "use strict";
    angular.module('mewpipe')
        .controller('MainController', ['$router', '$scope', 'userFactory', 'notificationFactory', 'themesFactory','$cookies', MainController])
    ;

    function MainController($router, $scope, userFactory, notificationFactory, themesFactory, $cookies) {
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
            {path: '/update-user/:id', component: 'updateUser'},
            {path: '/user/:id', component: 'user'}
        ]);
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
}());