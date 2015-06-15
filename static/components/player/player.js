(function () {
    "use strict";
    angular.module('mewpipe.player', [
        "ngSanitize"
    ])
        .controller('PlayerController', ["$sce", '$routeParams', 'mediaFactory', '$location', PlayerController])
    ;

    function PlayerController($sce, $routeParams, mediaFactory, $location) {
        var me = this;
        this.videoUrl = $sce.trustAsResourceUrl("/rest/media/" + $routeParams.id + "/read");
        this.mediaFactory = mediaFactory;
        var mediaId = mediaFactory.getMediaId();
        if (!mediaId || $routeParams.id !== mediaId){
            mediaFactory.setMediaId($routeParams.id);
            mediaFactory.getMedia($routeParams.id).success(function (response) {
                mediaFactory.setCurrentMedia(response);
                me.link = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/" + response.id;
                me.media = response;
                me.media.shares = me.media.shares || 0;
                me.media.views = me.media.views || 0;
            });
            this.user = function(){
                $location.url('/user/' + me.media.user.id);
            };
        }
    }
}());