package chit_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/bobg/chit"
)

func ExampleSQL() {
	var (
		ctx = context.Background()
		db  *sql.DB // Obtain a database handle from somewhere...
	)

	type rowtype struct {
		name, ssn string
	}
	const query = "SELECT name, ssn FROM employees WHERE salary >= 60000"

	iter := chit.SQL[rowtype](ctx, db, query)

	for {
		x, ok, err := iter.Next()
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			break
		}
		fmt.Printf("Employee %s, SSN %s\n", x.name, x.ssn)
	}
}
