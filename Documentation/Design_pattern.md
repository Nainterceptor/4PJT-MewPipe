Design Pattern
==============

## Architecture

* To build a powerful and scalable application, we chose to build a 3-tiers architecture. We have a Business part, 
composed of the database and the API; and a Front part, which is a single angularJS application for the moment.

* In the future, this architecture will allow us to easily build a mobile application or even a new website, 
thanks to the API!

## API

* Our API has been made in Golang, a google language. We wanted a fast, new and low level language for the backend part, 
and Golang is one of the best and mature language to do that.

* We have Entities, who are the models.
* Routes-docs files are use for the routing part, plus they document the request.
* And there is controllers files, who consume entities to perform operations.

We have a lot a additional things like filters, fixtures and unit tests.

## Front

* The Front is an angular application. We use MVVM pattern, so in our case the model is the Business part, 
the View are the html files, and the ViewModel are the javascripts components files.