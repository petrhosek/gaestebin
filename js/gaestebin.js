angular.module('gaestebin', ['ngResource', 'ngRoute', 'ngSanitize'])
    .config(function($routeProvider, $locationProvider) {
        $routeProvider
          .when('/paste/:id', {
            templateUrl: '/partials/paste-view.html',
            controller: 'PasteViewCtrl'
          })
          .when('/', {
            templateUrl: '/partials/paste-create.html',
            controller: 'PasteCreateCtrl'
          })
        $locationProvider.html5Mode(true);
    })
    .factory('Paste', function($resource) {
        return $resource('/api/paste/:id', {id:'@id'});
    })
    .service('popupService', function($window) {
        this.showPopup = function(message) {
          return $window.confirm(message);
        };
    })
    .controller('PasteListCtrl', function($scope, Paste) {
        $scope.pastes = Paste.query();
        $scope.deletePaste = function(paste) {
          if (popupService.showPopup('Really delete this?')) {
            paste.$delete(function() {
              $location.path('/');
            });
          }
        };
    })
    .controller('PasteViewCtrl', function($scope, $route, $routeParams, $location, popupService, Paste) {
        $scope.paste = Paste.get({id:$routeParams.id});
        $scope.deletePaste = function(paste) {
          if (popupService.showPopup('Really delete this?')) {
            paste.$delete(function() {
              $location.path('/');
            });
          }
        };
    })
    .controller('PasteCreateCtrl', function($scope, $location, Paste) {
        $scope.paste = new Paste();
        $scope.newPaste = function() {
          $scope.paste.$save(function(data) {
            console.log(data);
            $location.path('/paste/' + data.id);
          });
        };
    })
    .controller('PasteEditCtrl', function($scope, $route, $routeParams, Paste) {
        $scope.updatePaste = function() {
          $scope.paste.$update(function() {
            $location.path('/paste' + data.id);
          });
        };
        $scope.loadPaste = function() {
          $scope.paste = Paste.get({id:$routeParams.id});
        };
        $scope.loadPaste();
    })
    /*.controller('PasteCtrl', function($scope, $sce, $location, Paste) {
        $scope.pasteId = $location.path().substring(1);
        $scope.baseUrl = $location.absUrl().replace($location.path(), "");

        var highlightPaste = function(data) {
            $scope.paste = data;
            var highlighted = hljs.highlightAuto(data.Content);
            $scope.paste.highlighted = $sce.trustAsHtml(highlighted.value);
        }

        $scope.resetPaste = function() {
            $scope.pasteId = undefined;
            $scope.paste = undefined;
            $scope.pasteContent = undefined;
            $scope.pasteTitle = undefined;
        };

        $scope.newPaste = function() {
            var newPaste = new Paste();
            newPaste.Content = $scope.pasteContent;
            newPaste.Title = $scope.pasteTitle;
            var highlighted = hljs.highlightAuto($scope.pasteContent);
            newPaste.Language = highlighted.language;
            Paste.create(newPaste, function(data) {
                highlightPaste(data);
                $location.path('/' + data.Id);
            });
        };

        $scope.deletePaste = function() {
            console.log($scope.paste);
            Paste.delete({pasteId: $scope.paste.Id}, function(data) {
                console.log("Delete completed")
                console.log(data)
                $scope.resetPaste();
            });
        };

        $scope.showForm = function() {
            return !($scope.pasteId || $scope.paste);
        }

        if ($scope.pasteId.length > 0 &&
            (!$scope.paste || $scope.paste.Id != $scope.pasteId)) {
            Paste.get({pasteId: $scope.pasteId}, function(data) {
                highlightPaste(data);
            }, function(response) {
                console.log(response);
                $scope.resetPaste();
            });
        }
    })*/
    .directive('paperInput', function() {
      return {
        restrict: 'E',
        require: 'ngModel',
        link: function(scope, element, attrs, ctrl) {
          scope.$watch(function() {
            if (ctrl.$dirty) {
              return ctrl.$invalid
            } else {
              return false
            }
          }, function(invalid) {
            element[0].invalid = invalid;
          });
          
          element.on('input', function() {
            scope.$apply(function() {
              ctrl.$setViewValue(element.prop('inputValue'));
            });
          });
          ctrl.$setViewValue(null);

          ctrl.$render = function() {
            element.prop('inputValue', ctrl.$viewValue);
          };
        }
      };
    });
