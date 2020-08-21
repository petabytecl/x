package fmtx

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

func JSONPrettySPrint(p interface{}) (string, error) {
	json, err := jsoniter.MarshalIndent(p, "", "    ")
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%s", json), nil
}

func JSONPrettyPrint(p interface{}) {
	json, err := JSONPrettySPrint(p)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	fmt.Printf("%s\n", json)
}
