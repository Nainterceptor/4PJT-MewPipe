Mewpipe
=======

Mewpipe is the end-of-the-year project of the Master of Science of Supinfo 2015. 
MewPipe is a website that will allows users to host, share and play online videos.

Your can find our presentation on "gh-pages":https://nainterceptor.github.io/4PJT-MewPipe
Warning : This scenario and roles are pure fiction, and are not real.

There are scenarios that must be handled for answering completely to the client needs:

* An user must be able to upload a video, and obtain a link to share it
* An user must be able to edit a video information (Only the creator can achieve this operation)
* A user must be able to play a video
* A user must be able to control the video while playing (pause, fast-forward, rewind, ...)

## The team

Our developers for this project are :

* Robin Lebert (202429@supinfo.com)
* Gaël Demette (169174@supinfo.com)
* Alexandre Vaast (170768@supinfo.com)

But, our work is inseparable from the system part accomplished by :

* Mohamed Amir BEN SLAMIA (207910@supinfo.com)
* Arnaud Pierre BOYER (148282@supinfo.com)

## The mindset

Our aim for this project is to discover and have fun with new technologies. Some of dependencies for this project had 
to be corrected to meet our needs.

* "go-restful":https://github.com/emicklei/go-restful/pull/211
* "videogular":https://github.com/2fdevs/videogular/pull/211

## Our choices

Our strategy is to draw a reliable and scalable infrastructure that is able to take a amount of connections. 
We want to prevent a deny of service that could occur from a large success of our platform. 
In this way, we were looking for performance, maintainability and community sharing.

* The back-end HTTP server has been designed with [Golang](http://golang.org/), 
It's a recent and powerful language that offers some of packages from the community.
* The front-end is build with Angular.

## We have, one more thing..

Here are some things that we did outside the scope of the project:

* Swagger : We have an interface acting like a REST client. Gets documentation and behavior directly from the API !
* GoConvey : We have a 100% test coverage of our model layer. Tests are runnable from _./test.sh_ 
and gets a beautiful interface :)
* Cross-compilation : Our project can be compiled from linux to linux/Mac OS X and Windows. 
Just run _./build.sh_ and go to the _build_ directory
* Some fixtures : Our project can be filled by dummy data using the model layer. Just run _./fixtures.sh_
* A config file must be used to set parameters, each parameter must be overwritten directly from the CLI.
* In future releases, we have plan to make C.a.a.S (Coffee as a Service)

## Look 'n Feel

We have been inspired by the amazing work of [Aurélien Salomon](https://www.behance.net/aureliensalomon), 
and especially his [Youtube Redesign](https://dribbble.com/shots/1338727-Youtube-Redesign/attachments/189488). 
And with his [kind courtesy](https://twitter.com/aureliensalomon/status/609730220560617472). Thank you mate.

## Built in scripts

* Cross building with ./build.sh, outputs are in "build" directory
* Add some data with ./fixtures.sh
* Run the app on the fly with ./run.sh
* Run test suite with ./test.sh

## Documentation

* [Setup](Documentation/Install.md)
* [Database structure and optimization](Documentation/Database.md)
* [Ergonomics and ease to use](Documentation/Ergonomics.md)
* [Suggested Video algorithm explanation](Documentation/Suggested_algorithm.md)
* [Design Pattern](Documentation/Design_pattern.md)
