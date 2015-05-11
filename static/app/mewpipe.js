(function(){
    "use strict";
    angular.module('mewpipe',[
        'ngNewRouter',
        'mewpipeServices',
        'mewpipe.home',
        'mewpipe.player',
        'mewpipe.dashboard',
        'ui.bootstrap'
    ])
        .config(function($locationProvider){
            //$locationProvider.html5Mode(true);
        })
        .controller('MainController',['$router','$scope','userFactory','notificationFactory','themesFactory', MainController])
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

        this.logIn = function(){
            userFactory.logIn();
            notificationFactory.addAlert('Connected !','success');
        };
        this.logOut = function(){
            userFactory.logOut();
            notificationFactory.addAlert('Disconnected !','success');
        }
    }
}());