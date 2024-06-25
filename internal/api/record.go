package api

import (
	"DiscoDB/internal/models"
	"encoding/json"
)

// CreateRecord creates a new record in the specified table (text channel).
func (client *DBClient) CreateRecord(tableID string, record models.Record) error {
	// Marshal record to JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Send record as a message to the table (text channel)
	_, err = client.Session.ChannelMessageSend(tableID, string(recordJSON))
	if err != nil {
		return err
	}

	return nil
}

// ReadRecords reads all records from the specified table (text channel).
func (client *DBClient) ReadRecords(tableID string) ([]models.Record, error) {
	// Fetch messages (records) from the table (text channel)
	messages, err := client.Session.ChannelMessages(tableID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	var records []models.Record
	for _, message := range messages {
		var record models.Record
		err := json.Unmarshal([]byte(message.Content), &record)
		if err != nil {
			continue // Skip messages that aren't valid records
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord updates an existing record in the specified table (text channel).
func (client *DBClient) UpdateRecord(tableID, messageID string, record models.Record) error {
	// Marshal record to JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Edit the existing message (record) in the table (text channel)
	_, err = client.Session.ChannelMessageEdit(tableID, messageID, string(recordJSON))
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecord deletes an existing record from the specified table (text channel).
func (client *DBClient) DeleteRecord(tableID, messageID string) error {
	// Delete the message (record) from the table (text channel)
	err := client.Session.ChannelMessageDelete(tableID, messageID)
	if err != nil {
		return err
	}

	return nil
}
