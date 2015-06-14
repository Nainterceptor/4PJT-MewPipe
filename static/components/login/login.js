(function () {
    "use strict";
    angular.module('mewpipe.login', [
    ])
        .controller('LoginController', ['twitterFactory','$location',LoginController]);

    function LoginController(twitterFactory,$location) {
        var partUri = $location.path().split('/');
        var urlParams = $location.search();
        this.twitterConnect = function(){
            window.location.href = "/rest/users/login/twitter";
        };
        if(partUri[partUri.length - 1] === 'callback' && urlParams.oauth_token && urlParams.oauth_verifier){
            twitterFactory.login(urlParams.oauth_token, urlParams.oauth_verifier)
                .success(function(response){
                    console.log(response);
                })
                .error(function(response){
                    console.log(response);
                });
        }
    }
}());