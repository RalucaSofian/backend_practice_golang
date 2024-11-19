package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ExtractReqBody(req *http.Request, body any) error {
	reqBody, err := io.ReadAll(req.Body)
	fmt.Printf("[utils] err = %v; request body = %v\n", err, string(reqBody))
	if err != nil {
		fmt.Printf("[utils] err = %s", err.Error())
		return err
	}

	err = json.Unmarshal(reqBody, body)
	if err != nil {
		fmt.Printf("[utils] err = %s", err.Error())
		return err
	}

	return nil
}
