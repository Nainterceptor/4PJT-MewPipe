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
        .controller('MainController',['$router','userFactory','$scope', MainController]);

    function MainController($router, userFactory,$scope){
        $router.config([
            { path: '/', component: 'home'},
            { path: '/player', component: 'player'},
            { path: '/dashboard', component: 'dashboard'}
        ]);
        this.user = userFactory;

        this.logIn = function(){
            userFactory.logIn();
        };
        this.logOut = function(){
            userFactory.logOut();
        }
    }
}());