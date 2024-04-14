/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

// Package sqlschema includes schema definition and migration scripts for a SQL database.
package sqlschema

/*
SQL files places in this directory are executed during database initialization.
Each SQL script is executed only once, in order of the number in its file name.
Files must be named 1.sql, 2.sql, etc.
A script may contain multiple statements separated by a semicolon that is followed by a new line.

Typical schema definition and migration use cases:

	CREATE TABLE persons (
		tenant_id INT UNSIGNED NOT NULL,
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		name VARCHAR(256) CHARACTER SET ascii NOT NULL,
		created DATETIME(3) NOT NULL DEFAULT UTC_TIMESTAMP(3),
		PRIMARY KEY (tenant_id, person_id),
		INDEX(id),
		INDEX idx_name (name ASC)
	)

	ALTER TABLE persons
		DROP COLUMN created,
		ADD COLUMN first_name VARCHAR(128) CHARACTER SET ascii NOT NULL DEFAULT '',
		ADD COLUMN last_name VARCHAR(128) CHARACTER SET ascii NOT NULL DEFAULT ''

	CREATE INDEX idx_last_name ON persons (last_name ASC)
*/
