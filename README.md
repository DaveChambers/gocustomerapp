# gocustomerapp

A very simple web application with CRUD operations on a Customer object of the following form:

* First name - string (required, max length 100)
* Last name - string (required, max length 100)
* Birth date - datetime (required, from 18 till 60 years) 
* Gender - string (required, allowed values are Male, Female) 
*  E-mail - string (required, should be valid email)
* Address - string (optional, max length 200)

Written is Go and using HTML templates (html/template package) and Postgres inside Docker.  The app attempts to adhere to Uncle Bob's Clean Architecture and is influenced by Iman's excellent article:

https://medium.com/hackernoon/golang-clean-archithecture-efd6d7c43047

I thoroughly recommend Robert C Martin's wonderful book:

https://www.amazon.co.uk/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164

The app makes use of the following packages:

* github.com/brianvoe/gofakeit/v5
* github.com/bxcodec/faker
* github.com/gorilla/mux
* github.com/jinzhu/gorm
* github.com/joho/godotenv
* github.com/stretchr/testify

Integration tests are performed against a separate container and Unit tests use mocks.

## Docker Setup

Make 2 postgres instances inside Docker, with **separate** port numbers.  One is for the App and the second for the Integration Tests.

* Make a Docker postgres instance for the **app**:
``` bash
$ docker run --name customer-app -e POSTGRES_PASSWORD=yoursuperstr0ngpassword -d -p 5464:5432 postgres
```

* Make a Docker postgres instance for the **integration tests**:
``` bash
$ docker run --name customer-tests -e POSTGRES_PASSWORD=yoursuperstr0ngpassword -d -p 5465:5432 postgres
```

Copy `.env.example`, renaming it to `.env`, and enter your details.

Execute the `migrations/init.sql` against the both the App and Integration Test instances.

## Seeding

Seed the app instance via:

``` bash
$ go run seeder/seeder.go n
```

Where n = number of records.  

## Running

Run app with nodemon if you'd like to make changes without rerunning:

``` bash
$ nodemon --exec go run main.go --signal SIGTERM
```

Or just run:

``` bash
$ go run main.go
```

### Start here:

http://localhost:8080/create/

## Testing


Clean the cache: 

``` bash
$ go clean -testcache
```

Run all:

``` bash
$ go test ./...
```

Run all verbose:

``` bash
$ go test ./... -v
```