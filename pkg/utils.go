package pkg

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AdarshJha-1/Tempest/internal/config"
)

func PrettyPrintJSON(cfg config.Config) {
	prettyJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling: %s", err)
	}

	fmt.Println(string(prettyJSON))
}
