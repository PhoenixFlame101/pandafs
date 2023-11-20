package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a SQLite database. If the file doesn't exist, it will be created.
	db, err := sql.Open("sqlite3", "pandafs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the inode table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inode (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL,
			filesize INTEGER NOT NULL,
			isDirectory BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the chunk table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chunk (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			inode_id INTEGER,
			FOREIGN KEY (inode_id) REFERENCES inode(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the chunk_location table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chunk_location (
			chunk_id INTEGER,
			mule_id INTEGER,
			PRIMARY KEY (chunk_id, mule_id),
			FOREIGN KEY (chunk_id) REFERENCES chunk(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the directory table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS directory (
			file_id INTEGER PRIMARY KEY,
			dir_id INTEGER,
			FOREIGN KEY (file_id) REFERENCES inode(id),
			FOREIGN KEY (dir_id) REFERENCES inode(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database schema created successfully!")

	// Insert the root directory entry into the inode table
	rootDirFilename := "/"
	rootDirSize := 0
	isDirectory := true
	_, err = db.Exec(`
		INSERT INTO inode (id, filename, filesize, isDirectory) VALUES (0, ?, ?, ?)
	`, rootDirFilename, rootDirSize, isDirectory)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Root directory entry inserted successfully!")
}
