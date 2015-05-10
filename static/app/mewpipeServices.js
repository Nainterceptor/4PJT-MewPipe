(function(){
    "use strict";
    var config = {
        url: ""
    };

    angular.module('mewpipeServices',['ngCookies'])
        .factory('userFactory',['$http','$cookies', UserFactory])
        .factory('statsFactory',['$http', StatsFactory])
        .factory('notificationFactory',['$rootScope',NotificationFactory])
    ;

    function UserFactory($http, $cookies){
        var user = {};
        user.accessToken = $cookies.get('accessToken') ? $cookies.get('accessToken') : undefined;
        user.name = "Toto";
        user.email = "toto@tata.com";
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
                imgUrl: "http://lorempixel.com/400/300/nature/" + (i),
                title: "Title " + i,
                videoUrl: "/player"
            })
        }
        return {
            mostViewed: viewed,
            mostShared: viewed
        };
    }

    function NotificationFactory($rootScope){
        var factInstance = {};
        factInstance.alerts = [];

        factInstance.delAlert=function(index){
            factInstance.alerts.splice(index, 1);
        };

        factInstance.addAlert = function(msg, type, timer){
            console.log('alert');
            if(msg && type){
                factInstance.alerts.push({
                    type: 'alert-'+type,
                    msg: msg
                });
                delayDel(timer);

            }
        };

        var delayDel = function(timer){
            if(timer !== 0){
                var time = timer ? timer : 3000;
                setTimeout(function(){
                    factInstance.alerts.splice(-1 ,1);
                    $rootScope.$broadcast('alert:updated');
                }, time);
            }
        };

        return factInstance;
    }

}());