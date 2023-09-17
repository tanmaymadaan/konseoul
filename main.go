package main

import (
	"encoding/json"
	"fmt"
	"konseoul/common"
	"konseoul/db"
	"konseoul/targetgroup"
	"log"
	"net/http"
)

func main() {
	database := db.Connect()

	customMux := common.NewCustomServeMux()

	targetgroup.NewTargetGroupHandler(database, customMux)

	customMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := make(map[string]string)
		resp["message"] = "Healthy"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error happened in JSON marshal. Err: %s\n", err)
		}

		_, err = w.Write(jsonResp)
		if err != nil {
			log.Printf("Error happened in JSON marshal. Err: %s", err)
		}
	})

	registeredRoutes := customMux.GetRegisteredPatterns()
	fmt.Println("Registered routes")
	for i := 0; i < len(registeredRoutes); i++ {
		fmt.Println("\t ---", registeredRoutes[i])
	}

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", customMux)
	if err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
