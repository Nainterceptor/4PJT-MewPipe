(function () {
    "use strict";
    angular.module('mewpipe.updateVideo', [])
        .controller('UpdateVideoController', ['userFactory', 'notificationFactory', 'mediaFactory', '$timeout', '$location','$routeParams', UpdateVideoController])
        .config(function ($sceProvider) {
            $sceProvider.enabled(false);
        });

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

    function UpdateVideoController(userFactory, notificationFactory, mediaFactory, $timeout, $location,$routeParams) {
        var me = this;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        var currentMedia = mediaFactory.getCurrentMedia();
        if(!currentMedia || $routeParams.id != currentMedia.id){
            mediaFactory.getMedia($routeParams.id).success(function (response) {
                mediaFactory.setCurrentMedia(response);
                me.link = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/" + response.id;
                me.media = response;
            });
        } else {
            me.media = currentMedia;
        }
        this.user = userFactory;
        var fileToUpload;
        var thumbnail;
        me.canUpload = false;
        this.mediaFactory = mediaFactory;
        me.mediaScopes = ["public", "link", "private"];
        me.mediaScope = me.mediaScopes[0];
        me.prog = 0;
        this.uploading = false;
        this.validate = function (file) {
            me.canUpload = false;
            if (Object.getPrototypeOf(file).constructor === File) {
                var videoType = new RegExp("video/");
                if (!videoType.test(file.type)) {
                    notificationFactory.addAlert('Your file is not a video', 'danger');
                } else if (file.size > 524288000) {
                    notificationFactory.addAlert('Your file should not be superior to 500MB (file size : ' + ByteFilter(file.size) + ')', 'danger');
                }
                fileToUpload = file;
                var URL = window.URL;
                me.videoUrl = URL.createObjectURL(file);
                $timeout(function () {
                    var video = angular.element('#video')[0];
                    var canvas = angular.element('#canvas')[0];
                    canvas.width = 300;
                    canvas.height = 300 * video.videoHeight / video.videoWidth;
                    canvas.getContext('2d').drawImage(video, 0, 0, 300, 300 * video.videoHeight / video.videoWidth);
                    var img = canvas.toDataURL("image/png");
                    thumbnail = dataURItoBlob(img);
                    me.canUpload = true;
                    me.uploading = true;
                    me.prog = 0;
                    mediaFactory.upload(fileToUpload, thumbnail, me.media.id)
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
                }, 750);
            }
        };
        this.update = function(){
            mediaFactory.update(me.media).success(function(response){
                notificationFactory.addAlert('Media updated', 'success');
            });
        };
    }
}());