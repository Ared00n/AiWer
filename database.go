package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"  // Blank import для регистрации драйвера
)

var dbUsers *sql.DB
var dbCandidates *sql.DB
var dbWorks *sql.DB

func InitDB() {
	var err error

	// Создаём директорию db, если её нет
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		err = os.Mkdir("./db", 0755)
		if err != nil {
			log.Fatalf("Failed to create db directory: %v", err)
		}
	}

	// Инициализация users.db
	dbUsers, err = sql.Open("sqlite", "./db/users.db")
	if err != nil {
		log.Fatalf("Failed to open users.db: %v", err)
	}
	defer func() {  // Graceful close на ошибке
		if err != nil {
			dbUsers.Close()
		}
	}()

	_, err = dbUsers.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
        password TEXT
    )`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Инициализация works.db
	dbWorks, err = sql.Open("sqlite", "./db/works.db")
	if err != nil {
		log.Fatalf("Failed to open works.db: %v", err)
	}
	defer func() {
		if err != nil {
			dbWorks.Close()
		}
	}()

	_, err = dbWorks.Exec(`CREATE TABLE IF NOT EXISTS works (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        information TEXT NOT NULL,  -- Исправил 'informate' на 'information' (если не то — верни)
        start_date DATE NOT NULL,
        end_date DATE NOT NULL,
        time_duration INTEGER NOT NULL,
        collaborators INTEGER NOT NULL,
        username TEXT NOT NULL
        -- Убрал FOREIGN KEY: не работает cross-db. Свяжи по username в коде (JOIN) или используй одну БД
    )`)
	if err != nil {
		log.Fatalf("Failed to create works table: %v", err)
	}

	// Инициализация candidates.db
	dbCandidates, err = sql.Open("sqlite", "./db/candidates.db")
	if err != nil {
		log.Fatalf("Failed to open candidates.db: %v", err)
	}
	defer func() {
		if err != nil {
			dbCandidates.Close()
		}
	}()

	_, err = dbCandidates.Exec(`CREATE TABLE IF NOT EXISTS candidates (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        last_name TEXT,
        first_name TEXT,
        age INTEGER,
        profession TEXT,
        email TEXT UNIQUE,  -- Добавил UNIQUE для email (опционально, если нужно)
        module INTEGER,
        username TEXT
        -- Аналогично: если нужно FK на users, используй одну БД
    )`)
	if err != nil {
		log.Fatalf("Failed to create candidates table: %v", err)
	}

	// Пинг БД для проверки подключения
	if err = dbUsers.Ping(); err != nil {
		log.Fatalf("Users DB ping failed: %v", err)
	}
	if err = dbWorks.Ping(); err != nil {
		log.Fatalf("Works DB ping failed: %v", err)
	}
	if err = dbCandidates.Ping(); err != nil {
		log.Fatalf("Candidates DB ping failed: %v", err)
	}

	log.Println("All databases initialized successfully!")
}

func GetWorks() *sql.DB {
	return dbWorks
}

func GetDB() *sql.DB {
	return dbUsers
}

func GetCandidatesDB() *sql.DB {
	return dbCandidates
}

func CloseDB() {
	if dbUsers != nil {
		dbUsers.Close()
		log.Println("Users DB closed")
	}
	if dbCandidates != nil {
		dbCandidates.Close()
		log.Println("Candidates DB closed")
	}
	if dbWorks != nil {
		dbWorks.Close()
		log.Println("Works DB closed")
	}
}
