package pkg

import (
	"encoding/json"
	"fmt"
	"go/types"
	"log"
)

func PrettyPrintJSON(cfg types.Config) {
	prettyJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling: %s", err)
	}

	fmt.Println(string(prettyJSON))
}
