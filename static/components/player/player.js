(function () {
    "use strict";
    angular.module('mewpipe.player', [
        "ngSanitize",
        "com.2fdevs.videogular",
        "com.2fdevs.videogular.plugins.controls",
        "com.2fdevs.videogular.plugins.overlayplay",
        "com.2fdevs.videogular.plugins.poster"
    ])
        .controller('PlayerController', ["$sce", '$routeParams', 'mediaFactory', '$location', PlayerController])
    ;

    function PlayerController($sce, $routeParams, mediaFactory, $location) {
        var me = this;
        //this.config = {
        //    sources: [
        //        {src: $sce.trustAsResourceUrl("/rest/media/"+$routeParams.id+"/read"), type:"video/mp4"}
        //    ],
        //    theme: "bower_components/videogular-themes-default/videogular.css",
        //    plugins: {
        //    }
        //};
        this.videoUrl = $sce.trustAsResourceUrl("/rest/media/" + $routeParams.id + "/read");
        this.mediaFactory = mediaFactory;

        mediaFactory.getMedia($routeParams.id).success(function (response) {
            me.link = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/" + response.id;
            me.media = response;
        });
    }
}());