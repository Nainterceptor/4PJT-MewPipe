(function(){
    "use strict";
    var config = {
        url: ""
    };

    angular.module('mewpipeServices',['ngCookies'])
        .factory('userFactory',['$http','$cookies', UserFactory])
        .factory('statsFactory',['$http', StatsFactory])
    ;

    function UserFactory($http, $cookies){
        var user = {};
        user.accessToken = $cookies.get('accessToken') ? $cookies.get('accessToken') : undefined;
        user.logIn = function(){
            $cookies.put('accessToken','toto');
            user.accessToken = 'toto';
        };
        user.logOut = function(){
            $cookies.remove('accessToken');
            user.accessToken = undefined;
        };
        return user;
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