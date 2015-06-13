(function () {
    "use strict";
    var baseUrl = "/rest";

    angular.module('mewpipeServices', ['ngCookies', 'ngFileUpload'])
        .run(['$http', '$cookies', function ($http, $cookies) {
            if ($cookies.get('accessToken')) {
                $http.defaults.headers.common['Authorization'] = $cookies.get('accessToken');
            }
        }])
        .factory('userFactory', ['$http', '$cookies', 'notificationFactory', UserFactory])
        .factory('statsFactory', ['$http', StatsFactory])
        .factory('notificationFactory', ['$rootScope', NotificationFactory])
        .factory('themesFactory', ['$cookies', ThemesFactory])
        .factory('paginationFactory', [PaginationFactory])
        .factory('mediaFactory', ['$http', '$cookies', 'Upload', 'notificationFactory', MediaFactory])
    ;

    function UserFactory($http, $cookies, notificationFactory) {
        var userInstance = {};
        var isAdmin = function (user) {
            if (user.roles) {
                return user.roles.indexOf("Admin") !== -1;
            }
            return false;
        };
        userInstance.accessToken = $cookies.get('accessToken') ? $cookies.get('accessToken') : undefined;
        userInstance.getUsers = function () {
            return ($http.get(baseUrl + '/users'))
        };
        userInstance.setUser = function (user) {
            userInstance.user = user;
            userInstance.isAdmin = isAdmin(user);
        };
        userInstance.getUser = function () {
            if ((!userInstance.user || !userInstance.user.id) && $cookies.get('userId')) {
                userInstance.user = {
                    id: $cookies.get('userId')
                };
                return ($http.get(baseUrl + '/users/' + userInstance.user.id))
            }
            return false;
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
        userInstance.logIn = function (email, password) {
            $http.post(baseUrl + '/users/login', {
                email: email,
                password: password
            })
                .success(function (response) {
                    notificationFactory.addAlert('Connected !', 'success');
                    userInstance.user = response.User;
                    $cookies.put('accessToken', response.Token, {expires: new Date(response.ExpireAt)});
                    $http.defaults.headers.common['Authorization'] = response.Token;
                    $cookies.put('userId', response.User.id, {expires: new Date(response.ExpireAt)});
                    userInstance.accessToken = response.Token;
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
        userInstance.logOut = function () {
            $cookies.remove('accessToken');
            $cookies.remove('userId');
            userInstance.accessToken = undefined;
        };
        userInstance.updateUser = function (userId, email, firstname, lastname, nickname, password) {
            $http.put(baseUrl + '/users/' + userId, {
                email: email,
                name: {
                    firstname: firstname,
                    lastname: lastname,
                    nickname: nickname
                },
                password: password
            })
                .success(function (response) {
                    notificationFactory.addAlert('User updated !', 'success');
                })
                .error(function (response) {
                    notificationFactory.addAlert('Fail to update user', 'danger');
                    console.log(response);
                })
        };
        userInstance.deleteUser = function (userId) {
            $http.delete(baseUrl + '/users/' + userId)
                .success(function (response) {
                    notificationFactory.addAlert('User deleted !', 'danger');
                })
                .error(function (response) {
                    notificationFactory.addAlert('Fail to delete user', 'danger');
                    console.log(response);
                })
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
        mediaInstance.createMedia = function (user, title, summary) {
            return ($http.post(baseUrl + '/media', {
                title: title,
                user: user,
                summary: summary ? summary : ""
            }))
        };
        mediaInstance.getMedias = function () {
            return $http.get(baseUrl + '/media')
        };
        mediaInstance.getUserMedias = function () {
            return ($http.get(baseUrl + '/media/?user=' + $cookies.get('userId')))
        };
        mediaInstance.deleteMedia = function (mediaId) {
            $http.delete(baseUrl + '/media/' + mediaId, {
            })
                .success(function (response) {
                    notificationFactory.addAlert('Media deleted !', 'danger');
                })
                .error(function (response) {
                    notificationFactory.addAlert('Fail to delete media', 'danger');
                    console.log(response);
                })
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
        return mediaInstance;
    }

}());