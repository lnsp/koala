<p align="center">
<img src="https://user-images.githubusercontent.com/3391295/41180459-9f4527da-6b6e-11e8-8296-0979a1fc174b.png" alt="koala">
</p>
<hr>
koala is a simple web-service for editing local zonefiles. It HAS NO security built-in and should not be used outside any test environment.

## Dependencies
koala requires a recent Go version (tested on `>=1.10`) and npm.

## Configuration
```
KOALA_ADDR=":8080"
KOALA_ZONEFILE=
KOALA_STATICDIR="web/dist"
KOALA_APPLYCMD="sleep 1"
```
