(function () {
    "use strict";
    angular.module('mewpipe.dashboard', ['ngFileUpload'])
        .controller('DashboardController', ['userFactory', 'notificationFactory', '$location', DashboardController])
        .directive('profile', ['userFactory', 'notificationFactory', ProfileDirective])
        .directive('manageVideo', ['userFactory', 'notificationFactory', 'mediaFactory', 'paginationFactory', '$timeout', '$location', ManageVideoDirective])
        .config(function ($sceProvider) {
            $sceProvider.enabled(false);
        })
    ;

    function DashboardController(userFactory, notificationFactory, $location) {
        var me = this;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        this.user = userFactory;

        if ($location.url() != '/dashboard/manage-video') {
            this.activeTab = 'profile';
        } else {
            this.activeTab = 'video';
        }

        this.active = function (tab) {
            if (tab === 'profile') {
                $location.url('/dashboard/profile');
            } else {
                $location.url('/dashboard/manage-video');
            }
            me.activeTab = tab;
        };
    }

    function ProfileDirective(userFactory, notificationFactory) {
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/profile.html',
            scope: true,
            bindToController: true,
            controllerAs: 'profile',
            controller: function ($scope, $element, $attrs) {
                this.user = userFactory.user;
            }
        }
    }

    /**
     * @return {string}
     */
    function ByteFilter(bytes, precision) {
        if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
        if (typeof precision === 'undefined') precision = 1;
        var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
            number = Math.floor(Math.log(bytes) / Math.log(1024));
        return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) + ' ' + units[number];
    }

    function dataURItoBlob(dataURI) {
        // convert base64/URLEncoded data component to raw binary data held in a string
        var byteString;
        if (dataURI.split(',')[0].indexOf('base64') >= 0)
            byteString = atob(dataURI.split(',')[1]);
        else
            byteString = unescape(dataURI.split(',')[1]);

        // separate out the mime component
        var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0];

        // write the bytes of the string to a typed array
        var ia = new Uint8Array(byteString.length);
        for (var i = 0; i < byteString.length; i++) {
            ia[i] = byteString.charCodeAt(i);
        }

        return new Blob([ia], {type: mimeString});
    }

    function ManageVideoDirective(userFactory, notificationFactory, mediaFactory, paginationFactory, $timeout, $location) {
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/manage-video.html',
            scope: true,
            bindToController: true,
            controllerAs: 'mngVideo',
            controller: function ($scope, $element, $attrs) {
                var fileToUpload;
                var thumbnail;
                var me = this;
                this.mediaFactory = mediaFactory;
                mediaFactory.getUserMedias()
                    .success(function (response) {
                        me.baseUrl = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/";
                        console.log(response);
                        me.media = response;
                        paginationFactory.setPagination(me.media);
                        me.page = paginationFactory.getParams();
                    });
                me.prog = 0;
                this.uploading = false;
                this.validate = function (file) {
                    if (Object.getPrototypeOf(file).constructor === File) {
                        var videoType = new RegExp("video/");
                        if (!videoType.test(file.type)) {
                            notificationFactory.addAlert('Your file is not a video', 'danger');
                        } else if (file.size > 524288000) {
                            notificationFactory.addAlert('Your file should not be superior to 500MB (file size : ' + ByteFilter(file.size) + ')', 'danger');
                        }
                        fileToUpload = file;
                        var URL = window.URL;
                        me.title = file.name;
                        me.videoUrl = URL.createObjectURL(file);
                        $scope.$emit('videoRendered');
                        $scope.$digest();
                    }
                };
                $scope.$on('videoRendered', function (videoRenderedEvent) {
                    $timeout(function () {
                        meta();
                    }, 0, false);
                });
                var meta = function () {
                    var video = angular.element('#video')[0];
                    var canvas = angular.element('#canvas')[0];
                    canvas.width = 300;
                    canvas.height = 300 * video.videoHeight / video.videoWidth;
                    canvas.getContext('2d').drawImage(video, 0, 0, 300, 300 * video.videoHeight / video.videoWidth);
                    var img = canvas.toDataURL("image/png");
                    thumbnail = dataURItoBlob(img);
                    console.log(me.img);
                };
                this.upload = function () {
                    console.log(me.title, me.summary);
                    console.log(fileToUpload);
                    me.uploading = false;
                    me.prog = 0;
                    mediaFactory.createMedia(userFactory.user, me.title, me.summary)
                        .success(function (response) {
                            console.log(response);
                            me.uploading = true;
                            mediaFactory.upload(fileToUpload, thumbnail, response.id)
                                .progress(function (evt) {
                                    var prog = parseInt(100.0 * evt.loaded / evt.total);
                                    console.log('progress: ' + prog + '% file :' + evt.config.file.name);
                                    me.prog = prog;
                                })
                                .success(function (response) {
                                    console.log(response);
                                    me.playerUrl = response.id;
                                    me.link = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/" + response.id;
                                })
                                .error(function (response) {
                                    console.log('fail', response);
                                });
                        })
                        .error(function (response) {
                            console.log('fail', response);
                        });
                }
            }
        }
    }
}());