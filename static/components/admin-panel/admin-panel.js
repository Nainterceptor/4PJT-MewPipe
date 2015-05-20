(function () {
    "use strict";
    angular.module('mewpipe.adminPanel', [])
        .controller('AdminPanelController', ['userFactory','paginationFactory', AdminPanelController])
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
        this.admin = {
            users: [
                {username: "Robert1", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert2",
                    email: "test@",
                    createdAt: "10-11-2015"
                },
                {username: "Robert3", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert4",
                    email: "test@",
                    createdAt: "10-11-2015"
                },
                {username: "Robert5", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert6",
                    email: "test@",
                    createdAt: "10-11-2015"
                },
                {username: "Robert7", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert8",
                    email: "test@",
                    createdAt: "10-11-2015"
                },
                {username: "Robert9", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert10",
                    email: "test@",
                    createdAt: "10-11-2015"
                },
                {username: "Robert11", email: "test@", createdAt: "10-11-2015"}, {
                    username: "Robert12",
                    email: "test@",
                    createdAt: "10-11-2015"
                }
            ]
        };

        paginationFactory.setPagination(me.admin.users, 0, 5);
    }
}());