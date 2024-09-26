package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	BucketName       string
	BucketFolderName string
	DBUrl            string
	DBProductionUrl  string
}

func (e *EnvVariables) GetEnvVariables() {
	error := godotenv.Load(".env")
	if error != nil {
		log.Fatal(error)
	}

	dburl := os.Getenv("PRODUCTIONDB")
	dburlProd := os.Getenv("DATABASEURL")

	e.DBProductionUrl = dburlProd
	e.DBUrl = dburl
}
