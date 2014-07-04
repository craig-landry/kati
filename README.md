kati
=========

Kati is a simple reverse proxy for use in software development.

An example use case for it might be as follows:
You have an AngularJS client running via grunt/node locally on port 9000 where it needs to connect to that connects to a backend web server running on another machine on port 8080.  You can start up kati on port 80 and have it direct requests to the different servers based on url.

### Usage
    kati --http-port 80 --proxy "/api/.* -> api.example.com:8080" --proxy "/ng/.* -> angular-app.example.com"
This tells kati to listen on port 80 and when requests come in matching the regular expressions, forward them on to the host and port paired with the expression.


Version
----

1.0



License
----

GPL v2