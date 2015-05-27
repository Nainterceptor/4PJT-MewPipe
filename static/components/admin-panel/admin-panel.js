(function () {
    "use strict";
    angular.module('mewpipe.adminPanel', [])
        .controller('AdminPanelController', ['userFactory', 'paginationFactory', AdminPanelController])
        .filter('startFrom', AdminPanelFilter)
    ;

    function AdminPanelFilter() {
        return function (input, start) {
            if (!input || !input.length) {
                return;
            }
            start = +start; //parse to int
            return input.slice(start);
        }
    }

    function AdminPanelController(userFactory, paginationFactory) {
        var me = this;
        this.user = userFactory;

        paginationFactory.setPagination(me.user.users, 0, 5);
    }
}());