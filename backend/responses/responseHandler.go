package responses

import (
	"encoding/json"
	"strconv"
)

func CreateResponse(code GameStatus, message string, groupID string) (string, error) {
	response := map[string]string{
		"code":    strconv.Itoa(int(code)),
		"info":    message,
		"groupID": groupID,
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}
