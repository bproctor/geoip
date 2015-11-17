# GeoIP tools
Some tools for working with MaxMind GeoIP2Lite.
https://www.maxmind.com/en/geoip2-databases

## create-blocks.go
Converts the blocks file from CIDR into start and end IP addresses

An example database table in MySQL

```
CREATE TABLE `blocks` (
  `start_ip`                       INT(10)       UNSIGNED NOT NULL,
  `end_ip`                         INT(10)       UNSIGNED NOT NULL,
  `geoname_id`                     INT(10)       UNSIGNED          DEFAULT NULL,
  `registered_country_geoname_id`  INT(10)       UNSIGNED          DEFAULT NULL,
  `represented_country_geoname_id` INT(10)       UNSIGNED          DEFAULT NULL,
  `is_anonymous_proxy`             TINYINT(1)    UNSIGNED NOT NULL DEFAULT '0',
  `is_satellite_provider`          TINYINT(1)    UNSIGNED NOT NULL DEFAULT '0',
  `postal_code`                    VARCHAR(20)            NOT NULL DEFAULT '',
  `latitude`                       DECIMAL(18,9)                   DEFAULT NULL,
  `longitude`                      DECIMAL(18,9)                   DEFAULT NULL,
  PRIMARY KEY (`start_ip`,`end_ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

Example for running create-blocks.go
```
$ go run create-blocks.go GeoLite2-City-CSV_20151103/GeoLite2-City-Blocks-IPv4.csv > blocks.csv
```


## create-locations.go
Takes the timezones out and creates a new table out of them, replacing the string
with the primary key to a timezone lookup table.  It will output both CSVs to
stdout.

Example database tables in MySQL:

```
CREATE TABLE `locations` (
  `geoname_id`              INT(10)       UNSIGNED NOT NULL DEFAULT '0',
  `continent_code`          VARCHAR(2)             NOT NULL DEFAULT '',
  `continent_name`          VARCHAR(40)            NOT NULL DEFAULT '',
  `country_iso_code`        VARCHAR(2)             NOT NULL DEFAULT '',
  `country_name`            VARCHAR(64)            NOT NULL DEFAULT '',
  `subdivision_1_iso_code`  VARCHAR(3)             NOT NULL DEFAULT '',
  `subdivision_1_name`      VARCHAR(100)           NOT NULL DEFAULT '',
  `subdivision_2_iso_code`  VARCHAR(3)             NOT NULL DEFAULT '',
  `subdivision_2_name`      VARCHAR(100)           NOT NULL DEFAULT '',
  `city_name`               VARCHAR(100)           NOT NULL DEFAULT '',
  `metro_code`              SMALLINT(5)   UNSIGNED DEFAULT NULL,
  `timezoneid`              SMALLINT(3)   UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`geoname_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

```
CREATE TABLE `timezones` (
  `timezoneid`  SMALLINT(3)  UNSIGNED NOT NULL AUTO_INCREMENT,
  `timezone`    VARCHAR(40)           NOT NULL DEFAULT '',
  PRIMARY KEY (`timezoneid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

Example for runnning create-locations.go:
```
$ go run create-locations.go GeoLite2/GeoLite2-City-Locations-en.csv > locations.csv
```
