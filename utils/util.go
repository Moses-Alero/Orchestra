package utils

import (
	"encoding/json"
	"log"
	"orchestra/models"
	"os"
)

const orchestraInfoJSON string = ".orchestra.json"

func StoreOrchestraInfo(Orchestra *models.Cluster) error {
	bytes, err := json.MarshalIndent(Orchestra, "", " ")
	if err != nil {
		log.Fatal("Error", err)
	}
	if err := os.WriteFile(orchestraInfoJSON, bytes, 0666); err != nil {
		log.Fatal("Error: ", err)
	}

	return nil
}

func GetOrchestraInfo() *models.Cluster {
	bytes, err := os.ReadFile(orchestraInfoJSON)
	if err != nil {
		log.Fatal(err)
	}

	var orchestra models.Cluster
	err = json.Unmarshal(bytes, &orchestra)
	if err != nil {
		log.Fatal(err)
	}

	return &orchestra
}

func CheckForOrchestraInfo() bool {
	_, err := os.Stat(orchestraInfoJSON)
	return err == nil
}

func RemoveOrchestraInfo() {
	os.Remove(orchestraInfoJSON)
}
