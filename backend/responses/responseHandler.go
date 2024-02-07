package responses

import (
	"encoding/json"
	"strconv"
)

func CreateResponse(code GameStatus, message string, groupID string, data ...interface{}) (string, error) {
	response := map[string]interface{}{
		"code":    strconv.Itoa(int(code)),
		"info":    message,
		"groupID": groupID,
	}
	if len(data) > 0 {
		response["data"] = data[0]
	}
	jsonString, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}
