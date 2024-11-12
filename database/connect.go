package database

import (
	"database/sql"
	"fmt"
	"log"
	"myapi/config"
	"myapi/models"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Invalid port number")
		return
	}

	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		port,
	)

	dbName := config.Config("DB_NAME")

	sqlDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		panic("failed to connect to MySQL server")
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		panic("failed to create database: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		port,
		dbName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Video{},
		&models.VideoFormat{},
		&models.Token{},
		&models.Comment{},
	)
	if err != nil {
		log.Println("Error during migration:", err)
		return
	}

	fmt.Println("Database Migrated")
}
