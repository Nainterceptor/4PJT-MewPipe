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
        .controller('MainController',['$router','$scope','userFactory','notificationFactory', MainController])
    ;

    function MainController($router, $scope, userFactory,notificationFactory){
        var me = this;
        var themes = {
            default: {
                name: "default",
                url: "bootstrap/css/bootstrap.min.css"
            },
            slate: {
                name: "slate",
                url: "https://bootswatch.com/slate/bootstrap.min.css"
            },
            cosmo: {
                name: "cosmo",
                url: "https://bootswatch.com/cosmo/bootstrap.min.css"
            }
        };
        $router.config([
            { path: '/', component: 'home'},
            { path: '/player', component: 'player'},
            { path: '/dashboard', component: 'dashboard'}
        ]);
        this.chooseTheme = function(theme){
            me.theme = themes[theme];
            console.log(theme);
        };
        me.theme = {
            name: "default",
            url: "bootstrap/css/bootstrap.min.css"
        };
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