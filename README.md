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

## Run

```
go run main.go
```

- go to [http://localhost:8080](http://localhost:8080)
