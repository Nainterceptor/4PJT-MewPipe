(function(){
    "use strict";
    angular.module('mewpipe',[
        'ngNewRouter',
        'mewpipeServices',
        'mewpipe.home',
        'mewpipe.player',
        'mewpipe.dashboard',
        'mewpipe.adminPanel',
        'mewpipe.login',
        'mewpipe.upload'
    ])
        .config(function($locationProvider){
            $locationProvider.html5Mode(true);
        });

}());