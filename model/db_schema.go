package model

const (
	Schema = `
  CREATE TABLE "user"
	(
	   "user_id"     INTEGER PRIMARY KEY AUTOINCREMENT,
	   "username"    TEXT UNIQUE,
	   "display_name"TEXT,
	   "image_url"   TEXT,
	   "password"    TEXT,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT
	);
  
  CREATE TABLE "shop"
	(
	   "shop_id"     INTEGER PRIMARY KEY AUTOINCREMENT,
	   "user_id"     INTEGER,
	   "name"        TEXT,
	   "image_url"   TEXT,
	   "domain"      TEXT,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT
	);
  
  CREATE TABLE "product"
	(
	   "product_id"  INTEGER PRIMARY KEY AUTOINCREMENT,
	   "shop_id"     INTEGER,
	   "name"        TEXT,
	   "image_url"   TEXT,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT
	);
  
  CREATE TABLE "stock"
	(
	   "product_id"     INTEGER PRIMARY KEY AUTOINCREMENT,
	   "price_normal"   INTEGER,
	   "price_discount" INTEGER,
	   "total"          INTEGER,
	   "remaining"      INTEGER,
	   "status"         INTEGER,
	   "create_time"    TEXT,
	   "update_time"    TEXT
	);
  
  CREATE TABLE "cart"
	(
	   "cart_id"     INTEGER PRIMARY KEY AUTOINCREMENT,
	   "user_id"     INTEGER,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT
	);
  
  CREATE TABLE "cart_detail"
	(
	   "cart_detail_id" INTEGER PRIMARY KEY AUTOINCREMENT,
	   "cart_id"        INTEGER,
	   "product_id"     INTEGER,
	   "amount"         INTEGER,
	   "status"         INTEGER,
	   "create_time"    TEXT,
	   "update_time"    TEXT
	);
  
  CREATE TABLE "invoice"
	(
	   "invoice_id"  INTEGER PRIMARY KEY AUTOINCREMENT,
	   "user_id"     INTEGER,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT
	);
  
  CREATE TABLE "invoice_detail"
	(
	   "invoice_detail_id" INTEGER PRIMARY KEY AUTOINCREMENT,
	   "invoice_id"        INTEGER,
	   "product_id"        INTEGER,
	   "price_paid"        INTEGER,
	   "amount"            INTEGER,
	   "status"            INTEGER,
	   "create_time"       TEXT,
	   "update_time"       TEXT
	);  `
)
