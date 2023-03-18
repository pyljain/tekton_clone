package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func NewPostgres(connectionString string) (*Postgres, error) {

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS RUNNER (id serial PRIMARY KEY, name VARCHAR ( 50 ) UNIQUE NOT NULL)`,
	)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS LINK
			(id serial PRIMARY KEY, 
			RUNNER_ID INT,
			REPOSITORY VARCHAR (100), 
			CONSTRAINT fk_runner
			FOREIGN KEY(RUNNER_ID)
				REFERENCES RUNNER(id))`,
	)
	if err != nil {
		return nil, err
	}

	return &Postgres{database}, nil
}

func (postgres *Postgres) InsertRunner(ctx context.Context, name string) (int, error) {
	_, err := postgres.conn.Exec(`INSERT INTO RUNNER (name) VALUES ($1)`, name)
	if err != nil {
		return -1, err
	}

	row := postgres.conn.QueryRow("SELECT ID, NAME FROM RUNNER WHERE NAME = $1", name)

	var id int
	var n string

	err = row.Scan(&id, &n)
	if err != nil {
		return -1, err
	}

	return id, err
}

func (postgres *Postgres) CreateLink(ctx context.Context, runnerId int, repoName string) error {
	_, err := postgres.conn.Exec(`INSERT INTO LINK (RUNNER_ID, REPOSITORY) VALUES ($1, $2)`, runnerId, repoName)
	if err != nil {
		return err
	}

	return nil
}

func (postgres *Postgres) FindRunnerByRepository(ctx context.Context, repo string) (int, error) {
	row := postgres.conn.QueryRow("SELECT RUNNER_ID FROM LINK WHERE REPOSITORY = $1", repo)

	var runnerId int

	err := row.Scan(&runnerId)
	if err != nil {
		return -1, err
	}

	return runnerId, err
}
