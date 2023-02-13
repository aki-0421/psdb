<!---
Thanks for reading 💖!

## Commit Message Naming Convention

- F for Feature (Additions, Fixes, Ajustments of functionalities, etc.)
- T for Testing (New tests / specs, Test refactoring, etc.)
- R for Refactor (Adjustments of code structure, naming, typing, comments, etc.)
- D for Documentation (Documentation, README, etc.)
- S for Style (Styling, Storybook, Theme, Visual Design Adjustments, etc.)
- V for Version (Versioning, Dependencies, Licensing, etc.)
- C for Configuration (Building, Linting, CLI Tooling, etc.)

-->

# psdb

Go written sql driver (PlanetScale HTTP API)

# example

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/aki-0421/psdb"
    "github.com/aki-0421/psdb/pkg/config"
    "github.com/google/uuid"
)

func main() {
	c, err := config.New(host, user, pass, name)
	if err != nil {
		println(err)
	}

	db, err := sql.Open("planetscale", c.FormatDSN())
	if err != nil {
		println(err)
	}

	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `SELECT id,created_at FROM users;`)
	if err != nil {
		println(err)
	}

	for rows.Next() {
		var (
			id         uuid.UUID
			created_at int64
		)

		if err := rows.Scan(&id, &created_at); err != nil {
			panic(err)
		}

		fmt.Printf("%v %v\n", id, created_at)
	}
}
```