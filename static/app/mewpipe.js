(function(){
    "use strict";
    angular.module('mewpipe',[
        'ngNewRouter',
        'mewpipe.home',
        'mewpipe.player',
        'mewpipe.dashboard',
        'ui.bootstrap'
    ])
        .controller('MainController',['$router', MainController]);

    function MainController($router){
        $router.config([
            { path: '/', component: 'home'},
            { path: '/player', component: 'player'},
            { path: '/dashboard', component: 'dashboard'}
        ]);
    }
}());