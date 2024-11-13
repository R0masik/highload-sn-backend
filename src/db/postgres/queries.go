package postgres

const (
	// queries
	initQuery = `
CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    birth_date TIMESTAMP,
    sex CHARACTER VARYING(16) CHECK (sex IN (
		'male',
		'female'
	)),
    biography TEXT,
    city TEXT,
    password_hash TEXT
);

CREATE TABLE IF NOT EXISTS sessions(
    user_id uuid REFERENCES users(id),
    token TEXT
);
`

	insertQuery = "INSERT INTO %s(%s) VALUES %s"
	selectQuery = "SELECT %s FROM %s"
	updateQuery = "UPDATE %s SET %s WHERE id IN (%s)"
	deleteQuery = "DELETE FROM %s WHERE id IN (%s)"

	whereQuerySec = "WHERE %s"
)
