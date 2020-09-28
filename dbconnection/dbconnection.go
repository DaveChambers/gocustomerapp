package dbconnection

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/DaveChambers/gocustomerapp/testhelper"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

// Connect ...
func Connect() *gorm.DB {

	pathToDotEnv, testing := testhelper.GetRootPath() // We may or may not be testing.  If we are use TEST_PORT and path will need adjusting

	err := godotenv.Load(path.Join(pathToDotEnv, ".env"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	user := os.Getenv("DB_USER")
	var port string
	if testing {
		port = os.Getenv("TEST_PORT")
	} else {
		port = os.Getenv("PORT")
	}
	database := os.Getenv("DB")
	password := os.Getenv("PASSWORD")

	connectionString := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=disable password=%v", host, port, user, database, password)

	db, connectionError := gorm.Open("postgres", connectionString)

	if connectionError != nil {
		fmt.Println(connectionError)
		panic("Failed to connect to the database")
	}

	return db

}
