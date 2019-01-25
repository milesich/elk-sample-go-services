package store

// Plain SQL statements used by the services.
// Create/Delete tables.
const (
	SQLCreateSchema    = `CREATE SCHEMA IF NOT EXISTS elk`
	SQLCreateUserTable = `CREATE TABLE IF NOT EXISTS elk.users (
		id integer PRIMARY KEY,
		name text NOT NULL
	)`
	SQLCreateTaskTable = `CREATE TABLE IF NOT EXISTS elk.tasks (
		id integer PRIMARY KEY,
		user_id integer REFERENCES elk.users (id),
		name text NOT NULL,
		done boolean NOT NULL
	)`
	SQLDropSchema = `DROP SCHEMA elk CASCADE`
)

// Plain SQL statements used by the services.
// Insert/Update/Select.
const (
	SQLAddUser = `INSERT INTO elk.users (
		name
	)
	VALUES ($1)
	RETURNING id`
	SQLGetUser = `SELECT * FROM elk.users
		WHERE id = $1`
	SQLGetUsers = `SELECT * FROM elk.users`
	SQLAddTask  = `INSERT INTO elk.tasks (
		user_id,
		name,
		done
	)
	VALUES ($1, $2, false)
	RETURNING id`
	SQLUpdateTask = `UPDATE elk.tasks 
		SET name = $2, done = $3
		WHERE id = $1`
	SQLGetUserTasks = `SELECT * FROM elk.tasks
		WHERE user_id = $1`
)
