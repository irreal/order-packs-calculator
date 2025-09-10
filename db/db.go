package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/irreal/order-packs/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

// creates a new database connection and initialize the schema
func NewDB(dbPath string) (*DB, error) {

	// make sure dir exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	conn, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

	// if we just created the db, create schema and init
	if !dbExists {
		if err := db.initSchema(); err != nil {
			return nil, fmt.Errorf("failed to initialize schema: %w", err)
		}
		if err := db.seedData(); err != nil {
			return nil, fmt.Errorf("failed to seed data: %w", err)
		}
		log.Println("Database initialized with schema and sample data")
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS packs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		size INTEGER NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		requested_item_count INTEGER NOT NULL,
		shipped_item_count INTEGER NOT NULL,
		packs_json TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// seed db with sample data
func (db *DB) seedData() error {

	packSizes := []int{250, 500, 1000, 2000, 5000}
	for _, size := range packSizes {
		_, err := db.conn.Exec("INSERT INTO packs (size) VALUES (?)", size)
		if err != nil {
			return fmt.Errorf("failed to insert pack size %d: %w", size, err)
		}
	}

	samplePacks := map[models.Pack]int{
		250: 1,
	}
	packsJSON, err := json.Marshal(samplePacks)
	if err != nil {
		return fmt.Errorf("failed to marshal sample packs: %w", err)
	}

	_, err = db.conn.Exec(`
		INSERT INTO orders (requested_item_count, shipped_item_count, packs_json, status, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		100, 250, string(packsJSON), string(models.OrderStatusNew), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert sample order: %w", err)
	}

	return nil
}

// load all packs from db
func (db *DB) GetPacks() (models.Packs, error) {
	rows, err := db.conn.Query("SELECT size FROM packs ORDER BY size")
	if err != nil {
		return nil, fmt.Errorf("failed to query packs: %w", err)
	}
	defer rows.Close()

	var packs models.Packs
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("failed to scan pack: %w", err)
		}
		packs = append(packs, models.Pack(size))
	}

	return packs, nil
}

// replace all packs with new set
func (db *DB) SavePacks(packs models.Packs) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM packs")
	if err != nil {
		return fmt.Errorf("failed to delete existing packs: %w", err)
	}

	for _, pack := range packs {
		_, err = tx.Exec("INSERT INTO packs (size) VALUES (?)", int(pack))
		if err != nil {
			return fmt.Errorf("failed to insert pack size %d: %w", int(pack), err)
		}
	}

	return tx.Commit()
}

// add new order
func (db *DB) SaveOrder(order *models.Order) error {
	packsJSON, err := json.Marshal(order.Packs)
	if err != nil {
		return fmt.Errorf("failed to marshal packs: %w", err)
	}

	_, err = db.conn.Exec(`
		INSERT INTO orders (requested_item_count, shipped_item_count, packs_json, status, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		order.RequestedItemCount, order.ShippedItemCount, string(packsJSON), string(order.Status), order.CreatedAt)

	return err
}

// get data for web ui
func (db *DB) GetLast10Orders() ([]*models.Order, error) {
	rows, err := db.conn.Query(`
		SELECT requested_item_count, shipped_item_count, packs_json, status, created_at 
		FROM orders 
		ORDER BY created_at DESC 
		LIMIT 10`)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		var packsJSON string
		var statusStr string

		err := rows.Scan(&order.RequestedItemCount, &order.ShippedItemCount, &packsJSON, &statusStr, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal([]byte(packsJSON), &order.Packs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal packs: %w", err)
		}

		order.Status = models.OrderStatus(statusStr)
		orders = append(orders, &order)
	}

	return orders, nil
}
