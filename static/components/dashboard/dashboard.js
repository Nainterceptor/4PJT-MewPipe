(function(){
    "use strict";
    angular.module('mewpipe.dashboard', ['ngFileUpload'])
        .controller('DashboardController',['userFactory','notificationFactory','$location',DashboardController])
        .directive('profile',['userFactory','notificationFactory', ProfileDirective])
        .directive('manageVideo',['userFactory','notificationFactory','mediaFactory','paginationFactory','$timeout','$sce', ManageVideoDirective])
        .config(function($sceProvider){
            $sceProvider.enabled(false);
        })
    ;

    function DashboardController(userFactory, notificationFactory, $location){
        this.canActivate = function(){
            if (!userFactory.accessToken){
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        var me = this;
        this.user = userFactory;

        if($location.url() != '/dashboard/manage-video'){
            this.activeTab = 'profile';
        } else {
            this.activeTab = 'video';
        }

        this.active = function(tab){
            if(tab === 'profile'){
                $location.url('/dashboard/profile');
            } else {
                $location.url('/dashboard/manage-video');
            }
            me.activeTab = tab;
        };
    }

    function ProfileDirective(userFactory,notificationFactory){
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/profile.html',
            scope: true,
            bindToController: true,
            controllerAs: 'profile',
            controller: function($scope, $element, $attrs){
                this.user = userFactory.user;
            }
        }
    }

    /**
     * @return {string}
     */
    function ByteFilter(bytes, precision){
        if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
        if (typeof precision === 'undefined') precision = 1;
        var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
            number = Math.floor(Math.log(bytes) / Math.log(1024));
        return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
    }

    function ManageVideoDirective(userFactory,notificationFactory,mediaFactory,paginationFactory,$timeout,$sce){
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/manage-video.html',
            scope: true,
            bindToController: true,
            controllerAs: 'mngVideo',
            controller: function($scope, $element, $attrs){
                var fileToUpload;
                var me = this;
                this.media = mediaFactory.userMedias;
                paginationFactory.setPagination(me.media, 0, 10);
                me.prog = 0;
                this.uploading = false;
                this.validate = function(file){
                    if(Object.getPrototypeOf(file).constructor === File){
                        var videoType = new RegExp("video/");
                        if(!videoType.test(file.type)){
                            notificationFactory.addAlert('Your file is not a video', 'danger');
                        } else if (file.size > 524288000){
                            notificationFactory.addAlert('Your file should not be superior to 500MB (file size : ' + ByteFilter(file.size) + ')', 'danger');
                        }
                        console.log(file);
                        fileToUpload = file;
                         var URL = window.URL;
                        console.log(URL.createObjectURL(file));
                        $sce.trustAsResourceUrl(URL.createObjectURL(file));
                        console.log($sce.getTrustedResourceUrl());
                        me.videoUrl = URL.createObjectURL(file);
                        $scope.$digest();
                        //var canvas = angular.element('#canvas');
                        //var video = angular.element('#video');
                        $timeout(function(){
                            //var canvas = window.document.getElementById('canvas');
                            //var video = window.document.getElementById('video');
                            var canvas = angular.element('#canvas')[0];
                            var video = angular.element('#video')[0];
                            me.vidHeight =video.videoHeight;
                            me.vidWidth =video.videoWidth;
                            console.log(me.vidHeight);
                            $scope.$digest();
                            $timeout(function(){
                                    canvas.getContext('2d').drawImage(video, 0, 0,300, 300 * video.videoHeight/ video.videoWidth);
                                me.img = canvas.toDataURL("image/png");
                                console.log(me.img);
                            },500);
                        },100);

                        //var fileReader = new FileReader();
                        //fileReader.readAsDataURL(file);
                        //fileReader.onload = function(e) {
                        //    $timeout(function() {
                        //        console.log(e.target.result);
                        //        me.video = e.target.result;
                        //    });
                        //};

                        meta(file);
                    }
                };
                var meta = function(file){
                    me.title = file.name;

                    $scope.$digest();
                    angular.element('#metaModal').appendTo('body').modal('show');
                };
                this.upload = function(){
                    console.log(me.title, me.summary);
                    console.log(fileToUpload);
                    me.uploading = false;
                    me.prog = 0;
                    mediaFactory.createMedia(userFactory.user,me.title,me.summary)
                        .success(function(response){
                            console.log(response);
                            me.uploading = true;
                            mediaFactory.upload(fileToUpload, response.id)
                                .progress(function(evt) {
                                    var prog = parseInt(100.0 * evt.loaded / evt.total);
                                    console.log('progress: ' + prog + '% file :'+ evt.config.file.name);
                                    me.prog = prog;
                                })
                                .success(function(response){
                                    console.log(response);
                                    me.videoUrl = response.id;
                                })
                                .error(function(response){
                                    console.log('fail',response);
                                });
                        })
                        .error(function(response){
                            console.log('fail',response);
                        });
                }
            }
        }
    }
}());