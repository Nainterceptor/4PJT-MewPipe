(function () {
    "use strict";
    angular.module('mewpipe.updateUser', [])
        .controller('UpdateUserController', ['userFactory', 'notificationFactory', '$routeParams', UpdateUserController]);

    function UpdateUserController(userFactory,notificationFactory,$routeParams) {
        var me = this;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            if (!$routeParams.id) {
                notificationFactory.addAlert('Id required <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken && $routeParams.id;
        };
        if($routeParams.id){
            userFactory.getUserById($routeParams.id).success(function(response){
                me.user = response;
            });
        }
        this.update = function(){
            userFactory.updateUser(me.user)
                .success(function(response){
                    notificationFactory.addAlert('User Updated', 'success');
                })
                .error(function(response){
                })
        }

    }
}());