(function () {
    "use strict";
    angular.module('mewpipe.player', [
        "ngSanitize",
        "com.2fdevs.videogular",
        "com.2fdevs.videogular.plugins.controls",
        "com.2fdevs.videogular.plugins.overlayplay",
        "com.2fdevs.videogular.plugins.poster"
    ])
        .controller('PlayerController', ["$sce", '$routeParams', 'mediaFactory', PlayerController])
    ;

    function PlayerController($sce, $routeParams, mediaFactory) {
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

        mediaFactory.getMedia($routeParams.id).success(function (response) {
            me.media = response;
        });
    }
}());