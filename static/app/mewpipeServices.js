(function () {
    "use strict";
    var baseUrl = "/rest";

    angular.module('mewpipeServices', ['ngCookies', 'ngFileUpload'])
        .run(['$http', '$cookies', function ($http, $cookies) {
            if ($cookies.get('accessToken')) {
                $http.defaults.headers.common['Authorization'] = $cookies.get('accessToken');
            }
        }])
        .factory('userFactory', ['$http', '$cookies', 'notificationFactory','$location', UserFactory])
        .factory('statsFactory', ['$http', StatsFactory])
        .factory('notificationFactory', ['$rootScope', NotificationFactory])
        .factory('themesFactory', ['$cookies', ThemesFactory])
        .factory('paginationFactory', [PaginationFactory])
        .factory('mediaFactory', ['$http', '$cookies', 'Upload', 'notificationFactory', MediaFactory])
        .factory('twitterFactory', ['$http', TwitterFactory])
    ;

    function UserFactory($http, $cookies, notificationFactory, $location) {
        var userInstance = {};
        var isAdmin = function (user) {
            if (user.roles) {
                return user.roles.indexOf("Admin") !== -1;
            }
            return false;
        };
        var stayConnected = function(expire){

        };
        userInstance.accessToken = $cookies.get('accessToken') ? $cookies.get('accessToken') : undefined;
        userInstance.getUsers = function () {
            return ($http.get(baseUrl + '/users'))
        };
        userInstance.getOneUser = function (userId) {
            return ($http.get(baseUrl + '/users/' + userId))
        };
        userInstance.setUser = function (user) {
            userInstance.user = user;
            userInstance.isAdmin = isAdmin(user);
        };
        userInstance.getUser = function () {
            return ($http.get(baseUrl + '/users/' + userInstance.user.id))
        };
        userInstance.getUserById = function(id){
            return $http.get(baseUrl + '/users/' + id);
        };
        if ($cookies.get('userId')) {
            userInstance.user = {
                id: $cookies.get('userId')
            };
            userInstance.getUser()
                .success(function (response) {
                    userInstance.setUser(response);
                })
                .error(function (response) {
                    console.log(response);
                });
        }
        userInstance.initiate = function(response){
            userInstance.user = response.User;
            userInstance.isAdmin = isAdmin(response.User);
            $cookies.put('accessToken', response.Token, {expires: new Date(response.ExpireAt)});
            $http.defaults.headers.common['Authorization'] = response.Token;
            $cookies.put('userId', response.User.id, {expires: new Date(response.ExpireAt)});
            userInstance.accessToken = response.Token;
            notificationFactory.addAlert('Connected !', 'success');
        };
        userInstance.logIn = function (email, password) {
            $http.post(baseUrl + '/users/login', {
                email: email,
                password: password
            })
                .success(function (response) {
                    userInstance.initiate(response);
                })
                .error(function (response) {
                    notificationFactory.addAlert('Invalid Email or Password', 'danger');
                    console.log('failed');
                    console.log(response);
                });
        };
        userInstance.signUp = function (email, nickname, password) {
            $http.post(baseUrl + '/users', {
                email: email,
                name: {nickname: nickname},
                password: password
            })
                .success(function (response) {
                    notificationFactory.addAlert('Registered !', 'success');
                    userInstance.logIn(email, password);
                })
                .error(function (response) {
                    notificationFactory.addAlert('Disconnected !', 'success');
                    console.log('failed');
                    console.log(response);
                });
        };
        userInstance.updateUser = function(user){
            return(
                $http.put(baseUrl + '/users/' + user.id, user)
            )
        };
        userInstance.logOut = function () {
            $cookies.remove('accessToken');
            $cookies.remove('userId');
            userInstance.accessToken = undefined;
            userInstance.isAdmin = undefined;
            $location.url("/");
        };
        userInstance.deleteUser = function (userId) {
            return($http.delete(baseUrl + '/users/' + userId))
        };

        return userInstance;
    }

    function StatsFactory($http) {
        var viewed = [];
        for (var i = 0; i < 12; i++) {
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

    function NotificationFactory($rootScope) {
        var factInstance = {};
        factInstance.alerts = [];

        factInstance.delAlert = function (index) {
            factInstance.alerts.splice(index, 1);
        };

        factInstance.addAlert = function (msg, type, timer) {
            if (msg && type) {
                factInstance.alerts.push({
                    type: 'alert-' + type,
                    msg: msg
                });
                delayDel(timer);
            }
        };

        var delayDel = function (timer) {
            if (timer !== 0) {
                var time = timer ? timer : 3000;
                setTimeout(function () {
                    factInstance.alerts.splice(-1, 1);
                    $rootScope.$broadcast('alert:updated');
                }, time);
            }
        };

        return factInstance;
    }

    function ThemesFactory($cookies) {
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
        factory.saveTheme = function (theme) {
            var date = new Date();
            date.setFullYear(date.getFullYear() + 5);
            $cookies.put('theme', theme, {expires: date});
        };
        factory.getTheme = function () {
            return $cookies.get('theme');
        };
        return factory;
    }

    function PaginationFactory() {
        var page = {};
        page.setPagination = function (items) {
            page = {
                totalItems: items.length,
                currentPage: 0,
                numPerPage: 5
            };
        };
        page.getParams = function () {
            return {
                totalItems: page.totalItems,
                currentPage: page.currentPage,
                numPerPage: page.numPerPage,
                numberOfPages: Math.ceil(page.totalItems / page.numPerPage)
            };
        };

        return page;
    }

    function MediaFactory($http, $cookies, Upload, notificationFactory) {
        var mediaInstance = {};
        var currentMedia = undefined;
        var mediaId = undefined;
        mediaInstance.setCurrentMedia = function(media){
            currentMedia = media;
        };
        mediaInstance.getCurrentMedia = function(){
            return currentMedia;
        };
        mediaInstance.setMediaId = function(id){
            mediaId = id;
            setTimeout(function(){
                mediaId = undefined;
            },500);
        };
        mediaInstance.getMediaId = function(){
            return mediaId;
        };
        mediaInstance.createMedia = function (user, title, summary, scope) {
            console.log(scope);
            return ($http.post(baseUrl + '/media', {
                title: title,
                user: user,
                summary: summary ? summary : "",
                scope: scope
            }))
        };
        mediaInstance.mediaShare = function (mediaId) {
            return $http.post(baseUrl + '/media/' + mediaId + '/share', {
            })
                .success(function (response) {
                    notificationFactory.addAlert('Thanks for sharing', 'success');
                })
                .error(function (response) {
                    console.log(response);
                })
        };
        mediaInstance.getMedia = function (mediaId) {
            return $http.get(baseUrl + '/media/' + mediaId)
        };
        mediaInstance.getMedias = function () {
            return $http.get(baseUrl + '/media')
        };
        mediaInstance.getMediasByShares = function () {
            return $http.get(baseUrl + '/media?order=-shares')
        };
        mediaInstance.getMediasByViews = function () {
            return $http.get(baseUrl + '/media?order=-views')
        };
        mediaInstance.getUserMedias = function () {
            return ($http.get(baseUrl + '/media/?user=' + $cookies.get('userId')))
        };
        mediaInstance.getOneUserMedias = function (userId) {
            return ($http.get(baseUrl + '/media/?user=' + userId))
        };
        mediaInstance.deleteMedia = function (mediaId) {
           return $http.delete(baseUrl + '/media/' + mediaId)
        };
        mediaInstance.upload = function (file, thumbnail, mediaId) {
            return (
                Upload.upload({
                    url: baseUrl + "/media/" + mediaId + "/upload",
                    file: [file, thumbnail],
                    fileFormDataName: ['file', 'thumbnail']
                })
            )
        };
        mediaInstance.update = function(media){
            return $http.put(baseUrl + '/media/' + media.id, media);
        };
        return mediaInstance;
    }

    function TwitterFactory($http){
        console.log('toto');
        var twitterInstance = {};

        twitterInstance.login = function(oauth_token, oauth_verifier){
            return $http.get(baseUrl + '/users/login/twitter/callback?oauth_token='+oauth_token+'&oauth_verifier='+oauth_verifier);
        };

        return twitterInstance;
    }

}());