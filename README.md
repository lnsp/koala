# koala

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