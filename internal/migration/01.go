package migration

// Query01 ...
const Query01 = `
	CREATE TABLE scripts (
		id uuid PRIMARY KEY,
		script VARCHAR NOT NULL,
		description VARCHAR
	)
`
