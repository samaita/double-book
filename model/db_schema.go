package model

const (
	Schema = `
	CREATE TABLE "user"
	(
	   "user_id"     INTEGER NOT NULL UNIQUE,
	   "username"    TEXT UNIQUE,
	   "display_name"TEXT,
	   "image_url"   TEXT,
	   "password"    TEXT,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT,
	   PRIMARY KEY("user_id")
	);
  
  CREATE TABLE "shop"
	(
	   "shop_id"     INTEGER NOT NULL UNIQUE,
	   "user_id"     INTEGER,
	   "name"        TEXT,
	   "image_url"   TEXT,
	   "domain"      TEXT,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT,
	   PRIMARY KEY("shop_id")
	);
  
  CREATE TABLE "product"
	(
	   "product_id"  INTEGER NOT NULL UNIQUE,
	   "shop_id"     INTEGER,
	   "name"        TEXT,
	   "image_url"   TEXT,
	   "stock_id"    INTEGER,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT,
	   PRIMARY KEY("product_id")
	);
  
  CREATE TABLE "stock"
	(
	   "stock_id"       INTEGER NOT NULL UNIQUE,
	   "price_normal"   INTEGER,
	   "price_discount" INTEGER,
	   "total"          INTEGER,
	   "remaining"      INTEGER,
	   "status"         INTEGER,
	   "create_time"    TEXT,
	   "update_time"    TEXT,
	   PRIMARY KEY("stock_id")
	);
  
  CREATE TABLE "stock_detail"
	(
	   "stock_detail_id" INTEGER NOT NULL UNIQUE,
	   "stock_id"        INTEGER,
	   "user_id"         INTEGER,
	   "amount"          INTEGER,
	   "status"          INTEGER,
	   "create_time"     TEXT,
	   "update_time"     TEXT,
	   PRIMARY KEY("stock_detail_id")
	);
  
  CREATE TABLE "cart"
	(
	   "cart_id"     INTEGER NOT NULL UNIQUE,
	   "user_id"     INTEGER,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT,
	   PRIMARY KEY("cart_id")
	);
  
  CREATE TABLE "cart_detail"
	(
	   "cart_detail_id" INTEGER NOT NULL UNIQUE,
	   "cart_id"        INTEGER,
	   "product_id"     INTEGER,
	   "amount"         INTEGER,
	   "status"         INTEGER,
	   "create_time"    TEXT,
	   "update_time"    TEXT,
	   PRIMARY KEY("cart_detail_id")
	);
  
  CREATE TABLE "invoice"
	(
	   "invoice_id"  INTEGER NOT NULL UNIQUE,
	   "user_id"     INTEGER,
	   "status"      INTEGER,
	   "create_time" TEXT,
	   "update_time" TEXT,
	   PRIMARY KEY("invoice_id")
	);
  
  CREATE TABLE "invoice_detail"
	(
	   "invoice_detail_id" INTEGER NOT NULL UNIQUE,
	   "invoice_id"        INTEGER,
	   "product_id"        INTEGER,
	   "price_paid"        INTEGER,
	   "amount"            INTEGER,
	   "status"            INTEGER,
	   "create_time"       TEXT,
	   "update_time"       TEXT,
	   PRIMARY KEY("invoice_detail_id")
	);  `
)
