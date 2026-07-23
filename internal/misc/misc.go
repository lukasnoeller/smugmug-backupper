package misc

import (
	"encoding/json"
	"fmt"
)

func Parse(body []byte, object any) error {

	err := json.Unmarshal(body, object)
	if err != nil {
		fmt.Println("error ocurred during unmarshalling of response body: ", err)
		return err
	}
	return nil
}
func PrettyPrint(v any) {
	// json.MarshalIndent(data, prefix, indent)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error pretty printing:", err)
		return
	}
	fmt.Println(string(b))
}
