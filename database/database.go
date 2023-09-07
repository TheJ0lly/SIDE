package database

import (
	"fmt"
	"os"
)

type Database struct {
	database_dir string
}

// This function will initialize a database instance.
// If location exists, it will return the Database pointer, otherwise it will return nil.
func Initialize_Database(location string) *Database {
	temp, err := os.Stat(location)

	if os.IsNotExist(err) {
		fmt.Printf("Proposed database location does not exist! - \"%s\"\n", location)
		return nil
	}

	if !temp.IsDir() {
		fmt.Printf("Proposed database location is not a folder! - \"%s\"\n", location)
		return nil
	}

	return &Database{database_dir: location}
}

// This function is useless now and will do nothing.
func Rellocate_Database(location string) {

}

// This function will write to the database folder, where the key is supposed to be the block hash, and the value is the data inside the block.
// If the file is successfully created and written to, it will return nil, otherwise it will return the according error.
func (db *Database) Write_New_File_To_DB(key string) error {
	path_to_file := fmt.Sprintf("%s\\%s", db.database_dir, key)

	file, err := os.Create(path_to_file)

	if err != nil {
		return err
	}

	file.Close()
	return nil
}

// This function will update the data under a file.
// If the file exists and the data is updated, it will return nil, otherwise it will return the according error.
func (db *Database) Update_File_From_DB(key string, value []byte) error {
	file_location := fmt.Sprintf("%s\\%s", db.database_dir, key)

	_, err := os.Stat(file_location)

	if os.IsNotExist(err) {
		fmt.Printf("The file does not exist in the database! - \"%s\"\n", file_location)
		return err
	}

	os.Remove(file_location)

	err = os.WriteFile(file_location, value, 0666)

	if err != nil {
		return err
	}

	return nil

}

// This function will read from the database folder, and will return the bytes stored in the folder.
// If the exists and it contains data, it will return a slice of bytes, otherwise it will return nil.
func (db *Database) Read_File_From_DB(key string) ([]byte, error) {
	file_location := fmt.Sprintf("%s\\%s", db.database_dir, key)

	bytes, err := os.ReadFile(file_location)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// This function returns the directory in which we store the local database
func (db *Database) Get_Database_Dir() string {
	return db.database_dir
}
