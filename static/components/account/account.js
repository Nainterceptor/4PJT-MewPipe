(function () {
    "use strict";
    angular.module('mewpipe.account', [
    ])
        .controller('AccountController', ['userFactory','notificationFactory', AccountController])
    ;

    function AccountController(userFactory,notificationFactory) {
        console.log('toto');
        var me = this;
        this.canActivate = function () {
            if (!userFactory.accessToken) {
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        userFactory.getUser().success(function(response){
            userFactory.setUser(response);
            me.user = response;
        });
        this.update = function(){
            userFactory.updateUser(me.user)
                .success(function(response){
                    userFactory.setUser(response);
                    notificationFactory.addAlert('User Updated', 'success');
                })
                .error(function(response){
                    console.log(response);
                })
        }

    }
}());