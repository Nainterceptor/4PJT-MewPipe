(function(){
    "use strict";
    angular.module('mewpipe.home',[
        'mewpipeServices'
    ])
        .controller('HomeController',['statsFactory', HomeController]);

    function HomeController(statsFactory){
        var me = this;
        console.log(statsFactory);
        this.mostViewed = [];
        angular.forEach(statsFactory.mostViewed, function(stat, key){
            if(key % 6 === 0){
                me.mostViewed.push([]);
            }
            me.mostViewed[Math.floor(key / 6)].push(stat);
        });
        this.mostShared = [];
        angular.forEach(statsFactory.mostShared, function(stat, key){
            if(key % 6 === 0){
                me.mostShared.push([]);
            }
            me.mostShared[Math.floor(key / 6)].push(stat);
        });
    }
}());