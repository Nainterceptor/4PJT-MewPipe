(function () {
    "use strict";
    angular.module('mewpipe.login', [
    ])
        .controller('LoginController', ['twitterFactory','$location','userFactory',LoginController]);

    function LoginController(twitterFactory,$location,userFactory) {
        var partUri = $location.path().split('/');
        var urlParams = $location.search();
        this.twitterConnect = function(){
            window.location.href = "/rest/users/login/twitter";
        };
        if(partUri[partUri.length - 1] === 'callback' && urlParams.oauth_token && urlParams.oauth_verifier){
            twitterFactory.login(urlParams.oauth_token, urlParams.oauth_verifier)
                .success(function(response){
                    userFactory.initiate(response);
                })
                .error(function(response){
                });
        }
        this.logIn = function () {
            userFactory.logIn(this.signInEmail, this.signInPassword);
        };
        this.signUp = function () {
            userFactory.signUp(this.signUpEmail, this.signUpNickname, this.signUpPassword);
        };
    }
}());