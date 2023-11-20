package master

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func isValidFilename(filename, pwd string) bool {

	reservedCharacters := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}

	for _, char := range reservedCharacters {
		if strings.Contains(filename, char) {
			return false
		}
	}

	val, err := inodeExists(pwd, filename)
	if err != nil {
		log.Fatal(err)
	}

	if val {
		return false
	}

	return true
}

func inodeExists(pwd, filename string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM inode i
		JOIN directory d ON i.id = d.file_id
		WHERE i.filename = ? AND d.dir_id = (SELECT id FROM inode WHERE filename = ?);
	`

	var count int
	err := db.QueryRow(query, filename, pwd).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func HandleCommand(payload string) {

	InitializeDB()

	var data map[string]string
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// for key, value := range data {
	// 	fmt.Printf("Key: %s, Value: %v\n", key, value)
	// }

	switch data["cmd"] {
	case "ls":
		node.ListCommand(data["pwd"])
	case "touch":
		node.TouchCommand(data["filename"], data["pwd"])
		fmt.Println("got here lol")
	case "cp":
		node.CopyCommand(data["filename1"], data["filename2"], data["pwd"])
	case "mv":
		node.MoveCommand(data["filename1"], data["filename2"], data["pwd"])
	case "rm":
		node.RemoveCommand(data["filename"], data["pwd"])
	case "cd":
		node.DirExists(data["filename"], data["pwd"])
	case "upload":
		node.UploadCommand(data["filename"], data["filesize"], data["pwd"])
		// case "download":
		// 	node.TouchCommand(data["filename"], data["pwd"])
	}

	CloseDB()
}

func InitializeDB() {
	var err error
	db, err = sql.Open("sqlite3", "db/pandafs.db")
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetDirID(pwd string) (int, error) {
	dirs := strings.Split(pwd, "/")
	dirs = append([]string{"/"}, dirs...)

	var cleanDirs []string
	for _, dir := range dirs {
		if dir != "" {
			cleanDirs = append(cleanDirs, dir)
		}
	}

	var currentDirID int = 0

	for _, dir := range cleanDirs {
		var dirID int
		query := `
			SELECT id FROM inode
			WHERE filename = ?
				AND isDirectory = 1
				AND (id IN (SELECT file_id FROM directory WHERE dir_id = ?)
				OR id = 0)
		`
		err := db.QueryRow(query, dir, currentDirID).Scan(&dirID)
		if err != nil {
			return 0, fmt.Errorf("directory %s not found", dir)
		}
		currentDirID = dirID
	}

	return currentDirID, nil
}

func (n *MasterNode) TouchCommand(filename, pwd string) {
	if !isValidFilename(filename, pwd) {
		return
	}

	result, err := db.Exec(`
      INSERT INTO inode (filename, filesize, isDirectory) VALUES (?, 0, false);
   `, filename)
	if err != nil {
		log.Fatal(err)
	}

	inodeID, _ := result.LastInsertId()

	id, err := GetDirID(pwd)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
      INSERT INTO directory (file_id, dir_id) VALUES (?, ?);
   `, inodeID, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File created successfully.")
}

func (n *MasterNode) RemoveCommand(filename, pwd string) {
	dirID, err := GetDirID(pwd)
	if err != nil {
		log.Fatal(err)
	}

	query := `SELECT isDirectory FROM inode i
			JOIN directory d ON i.id = d.file_id
			WHERE i.filename = ?
			AND d.dir_id = ?`

	var isDir bool
	err = db.QueryRow(query, filename, dirID).Scan(&isDir)
	if err != nil {
		log.Fatal(err)
	}

	if !isDir {
		_, err = db.Exec("DELETE FROM inode WHERE filename = ?", filename)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("DELETE FROM directory WHERE file_id = ? AND dir_id = ?", filename, dirID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("File/directory removed successfully.")
	} else {
		var fileCount int
		err = db.QueryRow("SELECT COUNT(*) FROM directory WHERE dir_id = ?", dirID).Scan(&fileCount)
		if err != nil {
			log.Fatal(err)
		}
		if fileCount != 0 {
			fmt.Println("ERROR: Directory is not empty")
		} else {
			_, err = db.Exec("DELETE FROM inode WHERE filename = ?", filename)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("DELETE FROM directory WHERE file_id = ? AND dir_id = ?", filename, dirID)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("File/directory removed successfully.")
		}
	}

	fmt.Println("File/directory removed successfully.")
}

func (n *MasterNode) MoveCommand(srcFilename, destFilename, pwd string) {
	if !isValidFilename(filepath.Base(destFilename), pwd) {
		return
	}

	fmt.Println(filepath.Join(pwd, filepath.Dir(srcFilename)))
	srcDirID, err := GetDirID(filepath.Join(pwd, filepath.Dir(srcFilename)))
	if err != nil {
		log.Fatal(err)
	}

	var srcID int
	query := `SELECT d.file_id
		FROM directory d
		INNER JOIN inode i ON i.id = d.file_id
		WHERE d.dir_id = ? AND i.filename = ?;`
	err = db.QueryRow(query, srcDirID, srcFilename).Scan(&srcID)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("file id: %d\n", srcID)

	destDirPath := filepath.Dir(destFilename)

	destDirID, err := GetDirID(destDirPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("dir id: %d\n", destDirID)

	_, err = db.Exec("UPDATE directory SET dir_id = ? WHERE file_id = ? AND dir_id = ?", destDirID, srcID, pwd)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE inode SET filename = ? WHERE id = ?", filepath.Base(destFilename), srcID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File moved successfully.")
}

func (n *MasterNode) CopyCommand(srcFilename, destFilename, pwd string) {
	if !isValidFilename(filepath.Base(destFilename), pwd) {
		return
	}

	srcID, err := GetDirID(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	destDirPath := filepath.Dir(destFilename)

	destDirID, err := GetDirID(destDirPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO directory (file_id, dir_id) VALUES (?, ?)", srcID, destDirID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File copied successfully.")
}

func (n *MasterNode) ListCommand(pwd string) {
	dirID, err := GetDirID(pwd)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`
		SELECT i.filename, i.isDirectory
		FROM inode i
		INNER JOIN directory d ON i.id = d.file_id
		WHERE d.dir_id = ?;
	`, dirID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		var isDirectory bool
		err := rows.Scan(&filename, &isDirectory)
		if err != nil {
			log.Fatal(err)
		}

		if isDirectory {
			fmt.Println(filename + "/")
		} else {
			fmt.Println(filename)
		}
	}
}

func (n *MasterNode) UploadCommand(filename, filesize, pwd string) {
	size, _ := strconv.Atoi(filesize)
	AddFileToDB(filename, size, pwd)
}

func (n *MasterNode) DirExists(dir, pwd string) bool {
	_, err := GetDirID(filepath.Join(pwd, dir))
	if err != nil {
		return false
	} else {
		return true
	}
}

func AddFileToDB(filename string, filesize int, pwd string) error {
	dirID, err := GetDirID(pwd)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO inode (filename, filesize, isDirectory) VALUES (?, ?, false);
	`, filename, filesize)
	if err != nil {
		return err
	}

	inodeID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO directory (file_id, dir_id) VALUES (?, ?);
	`, inodeID, dirID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
