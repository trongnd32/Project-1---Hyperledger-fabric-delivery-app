'use strict'

var app = angular.module('application', []);

//Angular Controller

app.controller('appController', function ($scope, appFactory) {

	$("#success_trans").hide();
	$("#success_create").hide();
	$("#error_trans").hide();
	$("#error_query").hide();
	$("#success_query").hide();
	$("#success_update").hide();
	$("#error_update").hide();

	$scope.queryTrans = function () {
		var id = $scope.trans_id;

		appFactory.queryTrans(id, function(data){
			$scope.query_trans = data;

			if ($scope.query_trans == "Could not find transaction"){
				$("#error_query").show();
				$("#success_query").hide();
			} else{
				$("#error_query").hide();
				$("#success_query").show();
				
			}
		});
	},

	$scope.newTrans = function(){
		appFactory.newTrans($scope.trans, function(data){
			$("#success_create").show();
			$scope.create_trans = data;
		});
	},

	$scope.queryAllTrans = function(){
		appFactory.queryAllTrans(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_trans = array;
		});
	},

	$scope.updateState = function(){
		appFactory.updateState($scope.transUpdate, function(data){
			$scope.update_state = data;
			if ($scope.update_state == "Could not find transaction"){
				$("#error_update").show();
				$("#success_update").hide();
			} else{
				$("#error_update").hide();
				$("#success_update").show();
				
			}
		});
	}

})


app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllTrans = function(callback){

    	$http.get('/get_all_trans/').success(function(output){
			callback(output)
		});
	}

	factory.queryTrans = function(id, callback){
    	$http.get('/get_trans/'+id).success(function(output){
			callback(output)
		});
	}

	factory.newTrans = function(data, callback){

		data.key = data.id;
		var trans = data.key + "-" + data.id + "-" + data.paid + "-" + data.received;

    	$http.get('/add_trans/'+trans).success(function(output){
			$("#success_create").show();
			callback(output)
		});
	}

	factory.updateState = function(data, callback){
		data.key = data.id;
		var trans = data.key + "-" + data.id + "-" + data.paid + "-" + data.received;
		// $("#success_update").show();
    	$http.get('/update_state/'+trans).success(function(output){
			//$("#success_update").show();
			callback(output)
		});
	}

	return factory;
});