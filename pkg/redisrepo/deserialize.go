package redisrepo

import (
	"encoding/json"
	"log"

	"github.com/niluwats/gochat/model"
	"github.com/redis/go-redis/v9"
)

type Document struct {
	ID      string `json:"id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}

func Deserialize(res any) []Document {
	switch v := res.(type) {
	case []any:
		if len(v) > 1 {
			total := len(v) - 1
			var docs = make([]Document, 0, total/2)

			for i := 1; i <= total; i = i + 2 {
				arrOfValue := v[i+1].([]any)
				value := arrOfValue[len(arrOfValue)-1].(string)

				doc := Document{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}

				docs = append(docs, doc)
			}
			return docs
		}
	default:
		log.Printf("different response type other than []any. type: %T", res)
		return nil
	}
	return nil
}

func DeserializeChat(docs []Document) []model.Chat {
	chats := []model.Chat{}
	for _, doc := range docs {
		var c model.Chat
		json.Unmarshal(doc.Payload, &c)
		c.ID = doc.ID
		chats = append(chats, c)
	}

	return chats
}

func DeserializeContactList(contacts []redis.Z) []model.ContactList {
	contactList := make([]model.ContactList, 0, len(contacts))

	for _, contact := range contacts {
		contactList = append(contactList, model.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}
	return contactList
}
