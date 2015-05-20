(function(){
    "use strict";
    var baseUrl = "/rest";

    angular.module('mewpipeServices',['ngCookies'])
        .run(['$http','$cookies',function($http,$cookies) {
            if($cookies.get('accessToken')){
                $http.defaults.headers.common['Authorization'] = $cookies.get('accessToken');
            }
        }])
        .factory('userFactory',['$http','$cookies','notificationFactory', UserFactory])
        .factory('statsFactory',['$http', StatsFactory])
        .factory('notificationFactory',['$rootScope',NotificationFactory])
        .factory('themesFactory',['$cookies', ThemesFactory])
        .factory('paginationFactory',[PaginationFactory])
    ;

    function UserFactory($http, $cookies, notificationFactory){
        var userInstance = {};
        userInstance.accessToken = $cookies.get('accessToken') ? $cookies.get('accessToken') : undefined;
        userInstance.logIn = function(email, password){
            $http.post(baseUrl+'/users/login',{
                email: email,
                password: password
            })
                .success(function(response) {
                    notificationFactory.addAlert('Connected !','success');
                    userInstance.user = response.User;
                    $cookies.put('accessToken', response.Token, {expires: new Date(response.ExpireAt)});
                    userInstance.accessToken = response.Token;
                })
                .error(function(){
                    console.log('failed');
                    console.log(response);
                });
        };
        userInstance.signUp = function(email, nickname, password){
            $http.post(baseUrl+'/users',{
                email: email,
                name:{nickname: nickname},
                password: password
            })
                .success(function(response) {
                    notificationFactory.addAlert('Registered !','success');
                    userInstance.logIn(email, password);
                })
                .error(function(){
                    console.log('failed');
                    console.log(response);
                });
        };
        userInstance.logOut = function(){
            $cookies.remove('accessToken');
            userInstance.accessToken = undefined;
        };
        return userInstance;
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

    function ThemesFactory($cookies){
        var factory = {};
        factory.themes = {
            Default: {
                name: "Default",
                url: "bootstrap/css/bootstrap.min.css"
            },
            Slate: {
                name: "Slate",
                url: "https://bootswatch.com/slate/bootstrap.min.css"
            },
            Cosmo: {
                name: "Cosmo",
                url: "https://bootswatch.com/cosmo/bootstrap.min.css"
            },
            Darkly: {
                name: "Darkly",
                url: "https://bootswatch.com/darkly/bootstrap.min.css"
            },
            United: {
                name: "United",
                url: "https://bootswatch.com/united/bootstrap.min.css"
            }
        };
        factory.saveTheme = function(theme){
            var date = new Date();
            date.setFullYear(date.getFullYear() + 5);
            $cookies.put('theme', theme, {expires: date});
        };
        factory.getTheme = function(){
            return $cookies.get('theme');
        };
        return factory;
    }

    function PaginationFactory(){
        var page = {};
        page.setPagination = function (items, currentPage, numPerPage) {
            page = {
                totalItems: items.length,
                currentPage: currentPage,
                numPerPage: numPerPage
            };
        };
        page.getParams = function() {
            return {
                currentPage: page.currentPage,
                numPerPage: page.numPerPage
            };
        };
        page.numberOfPages = function () {
            return Math.ceil(page.totalItems / page.numPerPage);
        };
        return page;
    }

}());