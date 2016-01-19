//
// MaxMind GeoIP -> Convert text timezone locations to a separate lookup timezone table
//
// This takes in an locations file and spits out the same CSV but with the timezone
// column replaced with an index into the a timezone table rather than having the
// timezone location spelled out.
//
// The two tables are (in MySQL syntax):
//
// CREATE TABLE `locations` (
//   `geoname_id`              INT(10)       UNSIGNED NOT NULL DEFAULT '0',
//   `locale_code`             VARCHAR(2)             NOT NULL DEFAULT '',
//   `continent_code`          VARCHAR(2)             NOT NULL DEFAULT '',
//   `continent_name`          VARCHAR(40)            NOT NULL DEFAULT '',
//   `country_iso_code`        VARCHAR(2)             NOT NULL DEFAULT '',
//   `country_name`            VARCHAR(64)            NOT NULL DEFAULT '',
//   `subdivision_1_iso_code`  VARCHAR(3)             NOT NULL DEFAULT '',
//   `subdivision_1_name`      VARCHAR(100)           NOT NULL DEFAULT '',
//   `subdivision_2_iso_code`  VARCHAR(3)             NOT NULL DEFAULT '',
//   `subdivision_2_name`      VARCHAR(100)           NOT NULL DEFAULT '',
//   `city_name`               VARCHAR(100)           NOT NULL DEFAULT '',
//   `metro_code`              SMALLINT(5)   UNSIGNED DEFAULT NULL,
//   `timezoneid`              SMALLINT(3)   UNSIGNED NOT NULL DEFAULT '0',
//   PRIMARY KEY (`geoname_id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
// CREATE TABLE `timezone` (
//   `timezoneid`  SMALLINT(3)  UNSIGNED NOT NULL AUTO_INCREMENT,
//   `timezone`    VARCHAR(40)           NOT NULL DEFAULT '',
//   PRIMARY KEY (`timezoneid`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
// TODO: Output to files instead of outputing both tables it stdout, this will
//       make it easier to import later on.
//

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Keep track of all the timezones we've seen in a map
var timezones = map[string]int{}

func main() {

	inFile, _ := os.Open(os.Args[1])
	defer inFile.Close()
	reader := csv.NewReader(inFile)

	// Print the header line
	fmt.Println("geoname_id,locale_code,continent_code,continent_name,country_iso_code,country_name,subdivision_1_iso_code,subdivision_1_name,subdivision_2_iso_code,subdivision_2_name,city_name,metro_code,time_zone")

	lineCount := 0
	timezoneid := 1 // The next timezone, 0 is assumed to be UTC
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				return
			}
			break
		}

		lineCount++

		// The first line is the header line
		if lineCount == 1 {
			continue
		}

		// Get the timezone id for this timezone location
		tzid, ok := timezones[record[12]]
		if !ok {
			// Create the new timezone location in the map if it doesn't exist
			// and assign it a new timezoneid
			timezones[record[12]] = timezoneid
			tzid = timezoneid
			timezoneid++
		}

		fmt.Printf("%s,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%s,%d\n",
			record[0], record[1], record[2], record[3], record[4], record[5],
			record[6], record[7], record[8], record[9], record[10], record[11], tzid)

	}

	fmt.Println("\ntimezoneid,timezone")
	// Dump a CSV for the timezones
	for k, v := range timezones {
		fmt.Printf("%d,\"%s\"\n", v, k)
	}
}
