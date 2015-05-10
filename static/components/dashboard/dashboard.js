(function(){
    "use strict";
    angular.module('mewpipe.dashboard', [])
        .controller('DashboardController',['userFactory','notificationFactory','$routeParams',DashboardController])
        .directive('profile',['userFactory','notificationFactory', ProfileDirective])
        .directive('manageVideo',['userFactory','notificationFactory', ManageVideoDirective])
    ;

    function DashboardController(userFactory, notificationFactory, $routeParams){
        this.canActivate = function(){
            if (!userFactory.accessToken){
                notificationFactory.addAlert('You need to be connected, return to <a class="alert-link" href="/">Home</a>', 'danger', 3000);
            }
            return userFactory.accessToken;
        };
        var me = this;

        if(!$routeParams.dashboardChoice){
            this.activeTab = 'profile';
        }

        this.active = function(tab){
            me.activeTab = tab;
        };
    }

    function ProfileDirective(userFactory,notificationFactory){
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/profile.html',
            controller: function($scope, $element, $attrs){
                console.log('profile');
            }
        }
    }

    function ManageVideoDirective(userFactory,notificationFactory){
        return {
            restrict: 'E',
            templateUrl: 'components/dashboard/manage-video.html',
            controller: function($scope, $element, $attrs){
                console.log('video');
            }
        }
    }
}());