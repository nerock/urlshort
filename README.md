# _Nerock URL Shortener_

Simple URL shortener that allows to create, retrieve and delete shortened URLs

## How to run
```
go build -o urlshort cmd/main.go
./urlshort
```

## Documentation
The API documentation is available at `/docs` endpoint and can the file can be edited in `docs/swagger.json`

## Environment variables

|ENV VAR|SUMMARY|DEFAULT|
|-------|-------|-------|
|PORT|HTTP Server port|8080|
|DBCONN|Sqlite DB connection string|urlshort.db|
|DOMAIN|Domain where the app is deployed to build short URLs|localhost|