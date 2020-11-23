package postgres

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func (d *Database) beginTransaction() error {
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	d.tx = tx
	return nil
}

func (d *Database) endTransaction() error {
	if d.tx != nil {
		if err := d.tx.Commit(); err != nil {
			return fmt.Errorf("failed to end transaction: %w", err)
		}

		d.tx = nil
	}

	return nil
}

type transactionFunc func(t *sqlx.Tx) error

func (d *Database) WithTransaction(function transactionFunc, commit bool) error {
	if d.tx == nil {
		if err := d.beginTransaction(); err != nil {
			return err
		}
	}

	if err := function(d.tx); err != nil {
		if rollbackErr := d.tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%w: (and failed to rollback: %v)", err, rollbackErr)
		}
		return err
	}

	if commit {
		if err := d.endTransaction(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) SaveOrUpdate(chunk model.MapChunk, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.NamedExec(
			`INSERT INTO chunks (x, y, data, trees_count) VALUES (:x, :y, :data, :trees_count)
			   ON CONFLICT (x, y) DO UPDATE 
			   SET trees_count = :trees_count`,
			chunk)

		return err
	}, commit)
}

func (d *Database) GetMapChunk(x, y int64) (result model.MapChunk, err error) {
	err = d.db.Get(
		&result,
		"SELECT * FROM chunks WHERE x=$1 AND y=$2",
		x, y)
	return
}

func (d *Database) GetChatMessages(offset int, count int) (result []model.ChatMessage, err error) {
	err = d.db.Select(&result,
		"SELECT * FROM chatmessages ORDER BY message_id DESC OFFSET $1 LIMIT $2",
		offset, count)
	return
}

func (d *Database) AddChatMessage(message model.ChatMessage) (id int64, err error) {
	err = d.db.Get(&id,
		"INSERT INTO chatmessages (sender_name, text) VALUES ($1, $2) RETURNING message_id",
		message.SenderName, message.Text)
	return
}

func (d *Database) UpdateCharacter(character model.Character, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := d.db.NamedExec(
			`UPDATE characters SET 
                      name=:name, max_population=:max_population, current_population=:current_population
			   WHERE id=:id`, &character)
		return err
	}, commit)
}

func (d *Database) AddBuildingLocation(buildingLoc model.BuildingLocation, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.Exec(
			`INSERT INTO buildinglocations (building_id, owner_id, location)
				VALUES ($1, $2, $3)`, buildingLoc.BuildingID, buildingLoc.OwnerID, pg.Array(buildingLoc.Location))

		return err
	}, commit)
}

func (d *Database) GetBuildingLocation(location [3]float32) (result model.BuildingLocation, err error) {
	row := d.db.QueryRow("SELECT * FROM buildinglocations WHERE location=$1", pg.Array(location))

	var locationArr []float64
	err = row.Scan(&result.BuildingID, &result.OwnerID, pg.Array(&locationArr))
	if err == nil {
		result.Location[0] = float32(locationArr[0])
		result.Location[1] = float32(locationArr[1])
		result.Location[2] = float32(locationArr[2])
	}

	return
}

func (d *Database) GetBuildings() (result []model.Building, err error) {
	err = d.db.Select(&result, "SELECT * FROM buildings")
	return
}

func (d *Database) GetBuildingLocations() (result []model.BuildingLocation, err error) {
	rows, err := d.db.Query("SELECT * FROM buildinglocations")
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var location model.BuildingLocation
		var locationArr []float64
		err = rows.Scan(&location.BuildingID, &location.OwnerID, pg.Array(&locationArr))
		if err != nil {
			return
		}

		location.Location[0] = float32(locationArr[0])
		location.Location[1] = float32(locationArr[1])
		location.Location[2] = float32(locationArr[2])
		result = append(result, location)
	}

	err = rows.Err()
	return
}

func (d *Database) GetCharacters(accountID int) (result []model.Character, err error) {
	err = d.db.Select(&result,
		`SELECT c.* FROM accountcharacters as a
    INNER JOIN characters as c
        ON c.id = a.character_id
WHERE account_id = $1`, accountID)

	return
}

func (d *Database) GetAccount(login string) (result model.Account, err error) {
	err = d.db.Get(&result, "SELECT * from accounts WHERE login = $1", login)
	return
}

func (d *Database) GetCharacter(id int) (result model.Character, err error) {
	err = d.db.Get(&result, "SELECT * FROM characters WHERE id = $1", id)
	return
}

func (d *Database) AddCharacter(character model.Character, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := d.db.NamedExec("INSERT INTO characters VALUES (DEFAULT, :name, DEFAULT, DEFAULT)", character)
		return err
	}, commit)
}

func (d *Database) DeleteCharacter(id int, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := d.db.Exec("DELETE FROM characters WHERE id = $1", id)
		return err
	}, commit)
}

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	EnableSSL bool
}

func NewDatabase(config Config) (db.Database, error) {
	var sslMode string
	if config.EnableSSL {
		sslMode = "verify-full"
	} else {
		sslMode = "disable"
	}

	database, err := sqlx.Connect("postgres",
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
			config.DBName,
			config.User,
			config.Password,
			config.Host,
			config.Port,
			sslMode))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return &Database{
		db: database,
	}, nil
}
