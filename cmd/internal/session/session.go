package session

import (
	"encoding/json"
	"fmt"
	"os"
)

type Session struct {
	SessionName string `json:"sessionName"`
	UserName    string `json:"userName"`
	Host        string `json:"host"`
	Port        int    `json:"port,omitempty"`
	Password    string `json:"password,omitempty"`
}

func New(sessionName, userName, host, password string, port int) Session {
	return Session{
		SessionName: sessionName,
		UserName:    userName,
		Host:        host,
		Port:        port,
		Password:    password,
	}
}

func NewFromFile(fileName string) ([]Session, error) {
	file, err := os.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("session.newFromFile: error while reading file %s: %w", fileName, err)
	}

	var sessionArray []Session

	err = json.Unmarshal(file, &sessionArray)
	if err != nil {
		return nil, fmt.Errorf("session.newFromFile: error while parsing json: %w", err)
	}

	return sessionArray, nil
}

func SaveToFile(sessions []Session) error {
	workDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("session.SaveToFile: error while saving to file: %w", err)
	}

	fileDirectory := workDir + "/.essh"
	fileName := fileDirectory + "/session.json"

	text, err := json.MarshalIndent(sessions, "", "\t")
	if err != nil {
		return fmt.Errorf("session.SaveToFile: error while using MarshalIndent: %w", err)
	}

	if _, err := os.Stat(fileDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(fileDirectory, 0755); err != nil {
			return fmt.Errorf("session.SaveToFile: error while creating directory: %w", err)
		}
	}

	err = os.WriteFile(fileName, text, 0666)
	if err != nil {
		return fmt.Errorf("session.SaveToFile: error while writing to file: %w", err)
	}
	return nil
}
