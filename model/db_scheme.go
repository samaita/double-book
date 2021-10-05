package model

const (
	dbTableUser = `
	CREATE TABLE "user_credential" (
		"uid"	TEXT NOT NULL UNIQUE,
		"credential_type"	INTEGER DEFAULT 1,
		"credential_access"	TEXT UNIQUE,
		"credential_key"	TEXT,
		"credential_salt"	TEXT,
		"create_time"	TEXT,
		"update_time"	TEXT,
		PRIMARY KEY("uid")
	);`

	dbTableUserInfo = `
	CREATE TABLE "user_info" (
		"uid"	TEXT NOT NULL UNIQUE,
		"first_name"	TEXT,
		"last_name"	TEXT,
		"thumbnail"	TEXT,
		"create_time"	TEXT,
		"update_time"	TEXT,
		"status"	INTEGER,
		PRIMARY KEY("uid")
	);
	`
)
