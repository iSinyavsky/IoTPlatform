package auth

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"iot-project/tools"
)

func WriteFormLogin(tx *sql.Tx, actorid uint32, loginRecords []LoginRecord, user *UnconfirmedUserClaims, password string) error {
	for _, loginRecord := range loginRecords {
		if loginRecord.Type == FORM_LOGIN {
			return fmt.Errorf("UserIxistt")
		}
	}
	_, err := tx.Exec("insert into login (id, actorid, type, email, password, createdat) values (nextval('loginid'), $1, $2, $3, $4, now()) ",
		actorid,
		FORM_LOGIN,
		user.Email,
		HashPassword(password))
	if err != nil {
		return err
	}
	return nil
}

// Set optional invitedby field for row with id == actorid in table actor.
func SetInvitedBy(tx *sql.Tx, actorid uint32, invitedBy uint32) error {
	var result uint32
	err := tx.QueryRow("update actor set invitedby = $1 where id = $2 returning id", invitedBy, actorid).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

// Read all login records for a certain actorid.
func ReadLoginRecordsByActorId(actorid uint32) ([]LoginRecord, error) {
	rows, err := tools.DBS.Query("select id, actorid, type, extid, createdat, password, email from login where actorid=$1", actorid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readLoginRecords(rows)
}

// Read all login records for a certain extid (google or facebook user id).
func ReadLoginRecordsByExtId(extid string) ([]LoginRecord, error) {
	rows, err := tools.DBS.Query("select id, actorid, type, extid, createdat, password, email from login where extid=$1", extid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readLoginRecords(rows)
}

// Read all login records for a certain email.
func ReadLoginRecordsByEmail(email string) ([]LoginRecord, error) {
	email = strings.ToLower(email)
	rows, err := tools.DBS.Query("select id, actorid, type, extid, createdat, password, email from login where lower(email)=$1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readLoginRecords(rows)
}

func ReadActorRecordByEmail(email string) ActorRecord {
	actor := &ActorRecord{}
	tools.DBS.QueryRow("SELECT actor.id, actor.name, actor.avatarid FROM login LEFT JOIN actor ON login.actorid = actor.id WHERE lower(login.email) = $1", email).Scan(&actor.Id, &actor.Name, &actor.AvatarId)
	fmt.Println(actor)
	return *actor
}

func readLoginRecords(rows *sql.Rows) ([]LoginRecord, error) {
	loginRecords := make([]LoginRecord, 0)
	for rows.Next() {
		loginRecord := LoginRecord{}
		err := rows.Scan(
			&loginRecord.Id,
			&loginRecord.ActorId,
			&loginRecord.Type,
			&loginRecord.ExtId,
			&loginRecord.CreatedAt,
			&loginRecord.Password,
			&loginRecord.Email)
		if err != nil {
			return nil, err
		}
		loginRecords = append(loginRecords, loginRecord)
	}
	return loginRecords, nil
}

// LoginRecord represents a single row of login table.
type LoginRecord struct {
	Id        uint64
	ActorId   uint32
	Type      uint16         // Type of login record either db.FORM_LOGIN, db.GOOGLE_LOGIN or db.FACEBOOK_LOGIN
	ExtId     sql.NullString // Unique id specific to either google or facebook, null for form login
	Password  sql.NullString // Salted password hash, null for google and facebook login records
	CreatedAt time.Time
	Email     string
}

// Read a single row from actor table specified by primary key id.
func ReadActorById(id uint32) (*ActorRecord, error) {
	row := tools.DBS.QueryRow("select id, name, createdat, invitedby, lang, admin from actor where id=$1", id)
	actorRecord := ActorRecord{}
	err := row.Scan(
		&actorRecord.Id,
		&actorRecord.Name,
		&actorRecord.CreatedAt,
		&actorRecord.InvitedBy,
		&actorRecord.Lang,
		&actorRecord.Admin,
	)
	if err != nil {
		return nil, err
	}
	return &actorRecord, nil
}

// ActorRecord represents a single row of actor table.
type ActorRecord struct {
	Id        uint32
	Name      string
	CreatedAt time.Time
	InvitedBy sql.NullInt64
	Lang      string
	AvatarId  sql.NullInt64
	Admin     int
}

// curl -k -X POST -b  --data "@body.json" https://lost.report/api/load_tree
// en5d7117a9000000033417317fbab7e60748367266ca7dc37e4bf616ffbc8871bfcb138d2437fa1db7
