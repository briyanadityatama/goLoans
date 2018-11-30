// Starts restful application providing LMS (loans managemenet system) functionality
package main

import (
	"github.com/briyanadityatama/goLoans/lms/cola"
	"github.com/briyanadityatama/goLoans/lms/cola/infra/repo"
	"github.com/briyanadityatama/goLoans/rest"
)

func main() {
	lms := cola.New(repo.NewMemoryClientRepo())
	server := rest.NewLoansServer("localhost:8080", "http://localhost:8080", lms)
	server.Start()
}
