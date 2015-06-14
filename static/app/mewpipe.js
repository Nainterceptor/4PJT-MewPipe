(function(){
    "use strict";
    angular.module('mewpipe',[
        'ngNewRouter',
        'mewpipeServices',
        'mewpipe.home',
        'mewpipe.player',
        'mewpipe.adminPanel',
        'mewpipe.login',
        'mewpipe.upload',
        'mewpipe.account',
        'mewpipe.manageVideo',
        'mewpipe.updateVideo',
        'mewpipe.updateUser'
    ])
        .config(function($locationProvider){
            $locationProvider.html5Mode(true);
        });

}());