(function(){
    "use strict";
    var config = {
        url: ""
    };

    angular.module('mewpipeServices',[])
        .factory('userFactory',['$http', UserFactory])
        .factory('statsFactory',['$http', StatsFactory])
    ;

    function UserFactory($http){

    }

    function StatsFactory($http){
        var viewed = [];
        for(var i=0; i < 12; i++){
            viewed.push({
                imgUrl: "http://lorempixel.com/400/300",
                title: "Title " + i,
                videoUrl: "/player"
            })
        }
        return {
            mostViewed: viewed,
            mostShared: viewed
        };
    }
}());