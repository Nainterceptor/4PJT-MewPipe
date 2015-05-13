(function(){
    "use strict";
    angular.module('mewpipe')
        .controller('MainController',['$router','$scope','userFactory','notificationFactory','themesFactory', MainController])
        .controller('AuthenticationController',['userFactory','notificationFactory', AuthenticatitionController])
        .directive('modalAuth',[ModalAuthDirective])
    ;

    function MainController($router, $scope, userFactory,notificationFactory, themesFactory){
        var me = this;
        this.themes = themesFactory.themes;
        $router.config([
            { path: '/', component: 'home'},
            { path: '/player', component: 'player'},
            { path: '/dashboard', component: 'dashboard'}
        ]);
        this.chooseTheme = function(theme){
            me.theme = themesFactory.themes[theme];
            themesFactory.saveTheme(theme);
        };
        me.theme = themesFactory.getTheme()? themesFactory.themes[themesFactory.getTheme()] : themesFactory.themes['default'];
        this.alerts = notificationFactory.alerts;
        $scope.$on('alert:updated', function() {
            $scope.$apply();
        });
        this.closeAlert = function(index){
            notificationFactory.delAlert(index);
        };
        this.user = userFactory;
    }

    function AuthenticatitionController(userFactory, notificationFactory){
        this.test="toto";
        this.logIn = function(){
            userFactory.logIn();
            notificationFactory.addAlert('Connected !','success');
        };
        this.logOut = function(){
            userFactory.logOut();
            notificationFactory.addAlert('Disconnected !','success');
        }
    }

    function ModalAuthDirective(){
        return {
            restrict: 'E',
            templateUrl: 'app/templates/auth-modal.html',
            controller: function($scope){
                console.log($scope);
            }
        }
    }
}());