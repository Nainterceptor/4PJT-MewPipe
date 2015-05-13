(function(){
    "use strict";
    angular.module('mewpipe',[
        'ngNewRouter',
        'mewpipeServices',
        'mewpipe.home',
        'mewpipe.player',
        'mewpipe.dashboard',
        'mewpipe.adminPanel',
    ]);
        /*.config(function($locationProvider){
            //$locationProvider.html5Mode(true);
        })*/

}());