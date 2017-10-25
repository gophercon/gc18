package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	srv "github.com/gophercon/gc18/services/profile-srv/proto/record"
)

var (
	db            *sql.DB
	Url           = "root@tcp(127.0.0.1:3306)/profile"
	database      string
	profileSchema = `CREATE TABLE IF NOT EXISTS profiles (
id varchar(36) primary key,
name varchar(255),
owner varchar(255),
type integer,
display_name varchar(255),
blurb varchar(255),
url varchar(255),
location varchar(255),
created integer,
updated integer,
unique (name));`

	q = map[string]string{
		"delete": "DELETE from %s.%s where id = ?",
		"create": `INSERT into %s.%s (
				id, name, owner, type, display_name, blurb, url, location, created, updated) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"update":             "UPDATE %s.%s set name = ?, owner = ?, type = ?, display_name = ?, blurb = ?, url = ?, location = ?, updated = ? where id = ?",
		"read":               "SELECT * from %s.%s where id = ?",
		"list":               "SELECT * from %s.%s limit ? offset ?",
		"searchName":         "SELECT * from %s.%s where name = ? limit ? offset ?",
		"searchOwner":        "SELECT * from %s.%s where owner = ? limit ? offset ?",
		"searchNameAndOwner": "SELECT * from %s.%s where name = ? and owner = ? limit ? offset ?",
	}
	st = map[string]*sql.Stmt{}
)

func Init() {
	var d *sql.DB
	var err error

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		panic("Invalid database url")
	}

	if len(parts[1]) == 0 {
		panic("Invalid database name")
	}

	url := parts[0]
	rest := parts[1]

	opts := strings.Split(rest, "?")
	database = opts[0]
	dbopts := opts[1]
	fmt.Println("Connecting to ", url+"/"+"?"+dbopts)
	if d, err = sql.Open("mysql", url+"/"+"?"+dbopts); err != nil {
		log.Fatal(err)
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		log.Fatal(err)
	}
	d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(profileSchema); err != nil {
		log.Fatal(err)
	}
	db = d

	for query, statement := range q {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, "profiles"))
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}
}

func Create(profile *srv.Profile) error {
	profile.Created = time.Now().Unix()
	profile.Updated = time.Now().Unix()
	_, err := st["create"].Exec(profile.Id, profile.Name, profile.Owner, profile.Type, profile.DisplayName,
		profile.Blurb, profile.Url, profile.Location, profile.Created, profile.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(profile *srv.Profile) error {
	profile.Updated = time.Now().Unix()
	_, err := st["update"].Exec(profile.Name, profile.Owner, profile.Type, profile.DisplayName,
		profile.Blurb, profile.Url, profile.Location, profile.Updated, profile.Id)
	return err
}

func Read(id string) (*srv.Profile, error) {
	profile := &srv.Profile{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&profile.Id, &profile.Name, &profile.Owner, &profile.Type, &profile.DisplayName,
		&profile.Blurb, &profile.Url, &profile.Location,
		&profile.Created, &profile.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return profile, nil
}

func Search(name, owner string, limit, offset int64) ([]*srv.Profile, error) {
	var r *sql.Rows
	var err error

	if len(name) > 0 && len(owner) > 0 {
		r, err = st["searchNameAndOwner"].Query(name, owner, limit, offset)
	} else if len(name) > 0 {
		r, err = st["searchName"].Query(name, limit, offset)
	} else if len(owner) > 0 {
		r, err = st["searchOwner"].Query(owner, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var profiles []*srv.Profile

	for r.Next() {
		profile := &srv.Profile{}
		if err := r.Scan(&profile.Id, &profile.Name, &profile.Owner, &profile.Type, &profile.DisplayName,
			&profile.Blurb, &profile.Url, &profile.Location,
			&profile.Created, &profile.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		profiles = append(profiles, profile)

	}
	if r.Err() != nil {
		return nil, err
	}

	return profiles, nil
}
