package models

import (
	"database/sql"
	"fmt"
)

func InitSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TEXT,
		updated_at TEXT
	);

	CREATE TABLE IF NOT EXISTS folders (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TEXT,
		updated_at TEXT
	);

	CREATE TABLE IF NOT EXISTS vault_cids (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		cid TEXT NOT NULL,
		tx_hash TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS login_entries (
		id TEXT PRIMARY KEY,
		entry_name TEXT,
		folder_id TEXT,
		type TEXT,
		additionnal_note TEXT,
		custom_fields TEXT,
		created_at TEXT,
		updated_at TEXT,
		user_name TEXT,
		password TEXT,
		web_site TEXT
	);

	CREATE TABLE IF NOT EXISTS card_entries (
		id TEXT PRIMARY KEY,
		entry_name TEXT,
		folder_id TEXT,
		type TEXT,
		additionnal_note TEXT,
		custom_fields TEXT,
		created_at TEXT,
		updated_at TEXT,
		owner TEXT,
		number TEXT,
		expiration TEXT,
		cvc TEXT
	);

	CREATE TABLE IF NOT EXISTS identity_entries (
		id TEXT PRIMARY KEY,
		entry_name TEXT,
		folder_id TEXT,
		type TEXT,
		additionnal_note TEXT,
		custom_fields TEXT,
		created_at TEXT,
		updated_at TEXT,
		genre TEXT,
		firstname TEXT,
		second_firstname TEXT,
		lastname TEXT,
		username TEXT,
		company TEXT,
		social_security_number TEXT,
		ID_number TEXT,
		driver_license TEXT,
		mail TEXT,
		telephone TEXT,
		address_one TEXT,
		address_two TEXT,
		address_three TEXT,
		city TEXT,
		state TEXT,
		postal_code TEXT,
		country TEXT
	);

	CREATE TABLE IF NOT EXISTS note_entries (
		id TEXT PRIMARY KEY,
		entry_name TEXT,
		folder_id TEXT,
		type TEXT,
		additionnal_note TEXT,
		custom_fields TEXT,
		created_at TEXT,
		updated_at TEXT
	);

	CREATE TABLE IF NOT EXISTS sshkey_entries (
		id TEXT PRIMARY KEY,
		entry_name TEXT,
		folder_id TEXT,
		type TEXT,
		additionnal_note TEXT,
		custom_fields TEXT,
		created_at TEXT,
		updated_at TEXT,
		private_key TEXT,
		public_key TEXT,
		e_fingerprint TEXT
	);

	CREATE TABLE IF NOT EXISTS vault_contents (
		id TEXT PRIMARY KEY,
		user_id integer,
		cid TEXT,
		is_draft BOOLEAN,
		folders TEXT,
		entries TEXT,
		created_at TEXT,
		updated_at TEXT
	);
	
	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_vault_contents_user_id ON vault_contents(user_id);
	CREATE INDEX IF NOT EXISTS idx_vault_contents_cid ON vault_contents(cid);
	CREATE INDEX IF NOT EXISTS idx_vault_cids_user_id ON vault_cids(user_id);
	CREATE INDEX IF NOT EXISTS idx_vault_cids_cid ON vault_cids(cid);

	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("schema init error: %w", err)
	}

	fmt.Println("✅ DB schema fully initialized")
	return nil
}

