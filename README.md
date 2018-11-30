# goLoans

[![Go Report Card](https://goreportcard.com/badge/github.com/briyanadityatama/goLoans)](https://goreportcard.com/report/github.com/briyanadityatama/goLoans)

# Go Loans

Example financial **RESTful microservice** for taking Loans written entirely in **Go** (Golang). It was written using Test-Driven Development approach and follows the Domain-Driven Design guidelines.

## Features

- apply for a loan
  - possibility to take one loan per client
  - first loan up to 50000000
  - only 3 applications from one ip per day
- repay the loan - either partially or in full
- extend the loan of a given client

## Before RUN

Make sure this project placed on your GOPATH/go or whatever your GOPATH name.

## Run

```
go run main.go
```

- go to [http://localhost:8080](http://localhost:8080)

## POST & GET Method

Try to POST first data via POSTMAN

POST => `http://localhost:8080/clients`

Insert raw data in body e.g

```
{
	"ktpNumber" : "3522582509010002",
	"amount"    : 10000000,
	"term"      : 30,
	"birthDate" : "1 December 1994",
	"name"      : "Doe"
}
```

Try to GET the clients by ktpNumber

GET => `http://localhost:8080/clients/3522582509010002`

output will look like this

```
{
    "ktpNumber": "3522582509010002",
    "birthDate": "1 December 1994",
    "name": "Doe",
    "goLoans": {
        "links": [
            {
                "rel": "self",
                "href": "http://localhost:8080/clients/3522582509010002/goLoans"
            }
        ]
    }
}
```
