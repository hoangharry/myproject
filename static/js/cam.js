function ImgCtrl($scope, $http){
    $scope.imgs = [];
    $scope.working = false;

    var logError = function(data, status) {
        console.log('code '+status+': '+data);
        $scope.working = false;
    };

    var refresh = function() {
        return $http.get('/cam/').
          success(function(data) { $scope.imgs = data.Images; }).
          error(logError);
    };
    refresh().then(function() { $scope.working = false; });
}
let url = 'localhost:8080/cam'
let response = await fetch(url)
    .then(response => function(response){
        var res = response.json()
        document.getElementById()
    })