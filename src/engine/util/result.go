package util

import (
	"time"
	"regexp"
)

type Result struct {
	Code int `json:"code,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type ResultSimple struct {
	Code int `json:"code,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

type Message struct {
	Timestamp int64 `json:"time,omitempty"`
	Level string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

func SetMessage(level string, msg string) Message {
	t := time.Now().Unix()

	var message Message
	message.Timestamp = t
	message.Level = level
	message.Message = msg

	return message
}

func SetMessages(inputMessages []string) []Message {
	var messages []Message
	for _, msg := range inputMessages {
		re := regexp.MustCompile(`(\S+)\s+(.*)`)
		match := re.FindStringSubmatch(msg)
		/*
		if (isLevel == false) {
			//var message Message
			msg := fmt.Sprintf("Message level %s invalid! Must be INFO|ERROR|WARN|DEBUG|CMD", match[1])
			message := SetMessage("ERROR",msg)
			messages = append(messages, message)
		} else {
			if len(match) != 0 {
				message := SetMessage(match[1],match[2])
				messages = append(messages,message)
			}
		}
		//log.Println("tesHEREt",isLevel)
		*/
		if len(match) != 0 {
			message := SetMessage(match[1],match[2])
			messages = append(messages,message)
		}	
	}

	return messages
}

func SetResult(code int, messages []Message) Result {
	var result Result
	result.Code = code
	result.Messages = messages

	return result
}
