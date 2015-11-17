//
// MaxMind GeoIP CIDR -> start and end IP
//
// This command creates the blocks file that can be imported into the database
// It basically just converts the CIDR IP ranges into first and last IP in integer format
// then outputs the rest of the CSV like normal
//
// example:
// go run create-blocks.go GeoLite2-City-CSV_20151103/GeoLite2-City-Blocks-IPv4.csv > blocks.csv
//
// Then import into a table like this (in MySQL syntax):
// CREATE TABLE `blocks` (
//   `start_ip`                       INT(10)       UNSIGNED NOT NULL,
//   `end_ip`                         INT(10)       UNSIGNED NOT NULL,
//   `geoname_id`                     INT(10)       UNSIGNED          DEFAULT NULL,
//   `registered_country_geoname_id`  INT(10)       UNSIGNED          DEFAULT NULL,
//   `represented_country_geoname_id` INT(10)       UNSIGNED          DEFAULT NULL,
//   `is_anonymous_proxy`             TINYINT(1)    UNSIGNED NOT NULL DEFAULT '0',
//   `is_satellite_provider`          TINYINT(1)    UNSIGNED NOT NULL DEFAULT '0',
//   `postal_code`                    VARCHAR(20)            NOT NULL DEFAULT '',
//   `latitude`                       DECIMAL(18,9)                   DEFAULT NULL,
//   `longitude`                      DECIMAL(18,9)                   DEFAULT NULL,
//   PRIMARY KEY (`start_ip`,`end_ip`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
// Once imported and it looks good, swap the "blocks" table with this one
// to make it live.
//
package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func ip2long(ipAddr string) uint32 {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

var netmaskTable = map[uint8]uint32{
	0:  0,
	1:  2147483648,
	2:  3221225472,
	3:  3758096384,
	4:  4026531840,
	5:  4160749568,
	6:  4227858432,
	7:  4261412864,
	8:  4278190080,
	9:  4286578688,
	10: 4290772992,
	11: 4292870144,
	12: 4293918720,
	13: 4294443008,
	14: 4294705152,
	15: 4294836224,
	16: 4294901760,
	17: 4294934528,
	18: 4294950912,
	19: 4294959104,
	20: 4294963200,
	21: 4294965248,
	22: 4294966272,
	23: 4294966784,
	24: 4294967040,
	25: 4294967168,
	26: 4294967232,
	27: 4294967264,
	28: 4294967280,
	29: 4294967288,
	30: 4294967292,
	31: 4294967294,
	32: 4294967295,
}

func iprange(ipAddr string) (uint32, uint32) {
	p := strings.Split(ipAddr, "/")
	if len(p) < 2 {
		return 0, 0
	}
	nm, err := strconv.ParseUint(p[1], 10, 32)
	if err != nil {
		return 0, 0
	}
	netmask := netmaskTable[uint8(nm)]

	ip := ip2long(p[0])

	// MaxMind GeoIP includes the whatchmacalit (base address or whatever) and
	// broadcast addresses in the range otherwise to exclude these, the real
	// (usable) first ip address is firstip + 1 and the real last ip is lastip - 1
	firstip := (ip & netmask)
	lastip := (ip | ^netmask)
	return firstip, lastip
}

func main() {
	inFile, _ := os.Open(os.Args[1])
	defer inFile.Close()
	reader := csv.NewReader(inFile)

	// Print the header line
	fmt.Println("start_ip,end_ip,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider,postal_code,latitude,longitude")

	lineCount := 0
	for {
		// Read the line
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error: ", err.Error())
			}
			return
		}

		lineCount++

		// The first line is the header line
		if lineCount == 1 {
			continue
		}

		// Find the first and last IP
		first, last := iprange(record[0])

		// Write the line
		fmt.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n", first, last, record[1], record[2],
			record[3], record[4], record[5], record[6], record[7], record[8])
	}
}
