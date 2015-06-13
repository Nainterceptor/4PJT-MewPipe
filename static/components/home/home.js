(function(){
    "use strict";
    angular.module('mewpipe.home',[
        'mewpipeServices'
    ])
        .controller('HomeController',['mediaFactory', 'paginationFactory', '$location', HomeController]);

    function HomeController(mediaFactory, paginationFactory, $location){
        var me = this;
        this.mostShared = [];

        mediaFactory.getMedias().success(function (response) {
            me.baseUrl = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/";
            me.medias = response;
            paginationFactory.setPagination(me.medias);
            me.page = paginationFactory.getParams();
        });
    }
}());