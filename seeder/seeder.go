package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DaveChambers/gocustomerapp/dbconnection"
	"github.com/DaveChambers/gocustomerapp/domain"
	"github.com/brianvoe/gofakeit/v5"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func randomBirthday() time.Time {
	// A birthday within our allowed range will be generated
	min := time.Date(1961, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2001, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func randomAddress() string {
	street := gofakeit.StreetNumber() + " " + gofakeit.StreetName() + " Street"
	city := gofakeit.City()
	state := gofakeit.State()
	zip := gofakeit.Zip()
	country := gofakeit.Country()
	return street + ", " + city + ", " + state + ", " + zip + ", " + country
}

func main() {

	// Convert our command line input to an Int
	numberOfRecords, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	db := dbconnection.Connect()
	defer db.Close()

	gofakeit.Seed(0)

	dupeCount := 0

	for i := 0; i < numberOfRecords; i++ {

		customer := &domain.Customer{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			BirthDate: randomBirthday(),
			Gender:    gofakeit.Gender(),
			Email:     gofakeit.Email(),
			Address:   randomAddress()}

		// Handle seed data rejected by the database:
		if err := db.Create(&customer).Error; err != nil {

			fmt.Println(err.Error())

			// In the very unlikely event of gofakeit generating an email already in the db, try again:
			if err.Error() == "pq: duplicate key value violates unique constraint \"customers_email_key\"" {
				fmt.Println("Dupe email address => generate another...")
				numberOfRecords++ // Increment to allow an additional record to be created
				dupeCount++
			}

			// In the very unlikely event of gofakeit generating an Firstname, Lastname or Address that is too long, try again:
			if strings.Contains(err.Error(), "pq: value too long for type character varying") {
				fmt.Println("Firstname, Lastname or Address were too long => generate another...")
				numberOfRecords++ // Increment to allow an additional record to be created
				dupeCount++
			}
		}
	}

	fmt.Printf("Successfully seeded the database with %d records!\n", numberOfRecords-dupeCount)
}
