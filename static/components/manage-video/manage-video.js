(function () {
    "use strict";
    angular.module('mewpipe.manageVideo', [])
        .controller('ManageVideoController', ['userFactory', 'mediaFactory','notificationFactory', '$location','paginationFactory', ManageVideoController])
    ;

    function ManageVideoController(userFactory, mediaFactory, notificationFactory, $location,paginationFactory) {
        var me = this;
        this.mediaFactory = mediaFactory;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        mediaFactory.getUserMedias()
            .success(function (response) {
                me.baseUrl = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/player/";
                me.media = response;
                paginationFactory.setPagination(me.media);
                me.page = paginationFactory.getParams();
            });
        this.update = function(media){
            mediaFactory.setCurrentMedia(media);
            $location.url('/update-video/' + media.id);
        };
        this.delete = function(id){
            mediaFactory.deleteMedia(id)
                .success(function (response) {
                    notificationFactory.addAlert('Media deleted !', 'success');
                    mediaFactory.getUserMedias().success(function (response) {
                        me.media = response;
                        paginationFactory.setPagination(me.medias);
                        me.page = paginationFactory.getParams();
                    });
                })
                .error(function (response) {
                    notificationFactory.addAlert('Fail to delete media', 'danger');
                })
        }
    }
}());