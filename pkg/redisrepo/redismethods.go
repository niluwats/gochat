package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/niluwats/gochat/model"
	"github.com/redis/go-redis/v9"
)

func RegisterNewUser(username, password string) error {
	err := redisClient.Set(context.Background(), username, password, 0).Err()
	if err != nil {
		log.Println("error adding new user - ", err)
		return err
	}

	err = redisClient.SAdd(context.Background(), userSetKey(), username).Err()
	if err != nil {
		log.Println("error adding user in set - ", err)
		redisClient.Del(context.Background(), username)
		return err
	}
	return nil
}

func IsUserExists(username string) bool {
	return redisClient.SIsMember(context.Background(), userSetKey(), username).Val()
}

func IsUserAuthenticated(username, password string) error {
	pass := redisClient.Get(context.Background(), username).Val()
	if !strings.EqualFold(pass, password) {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}

func UpdateContactList(username, contact string) error {
	zs := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	err := redisClient.ZAdd(context.Background(), contactListKey(username), *zs).Err()

	if err != nil {
		log.Println("error while updating contact list. username: ", username, " contact: ", contact, err)
		return err
	}

	return nil
}

func CreateChat(c *model.Chat) (string, error) {
	chatKey := chatKey()
	by, _ := json.Marshal(c)

	res, err := redisClient.Do(context.Background(), "JSON.SET", chatKey, "$", string(by)).Result()
	if err != nil {
		log.Println("error setting chat json ", err)
	}

	log.Println("chat set successfully ", res)

	err = UpdateContactList(c.From, c.To)
	if err != nil {
		log.Println("error updating contact list of ", c.From, err)
	}

	err = UpdateContactList(c.To, c.From)
	if err != nil {
		log.Println("error updating contact list of ", c.To, err)
	}

	return chatKey, nil
}

func CreateFetchChatBetweenIndex() {
	res, err := redisClient.Do(context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]", username1, username2, username1, username2, fromTS, toTS)
	res, err := redisClient.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
	).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := Deserialize(res)

	chats := DeserializeChat(data)
	return chats, nil
}

func FetchContactList(username string) ([]model.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListKey(username),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}

	res, err := redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	if err != nil {
		log.Println("error fetching contact list.username: ", username, err)
		return nil, err
	}

	contactList := DeserializeContactList(res)
	return contactList, nil
}
