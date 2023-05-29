package redisrepo

import "context"

func RegisterNewUser(username, password string) error {
	err := redisClient.Set(context.Background(), username, password, 0).Err()
}

func IsUserAuthenticated() {

}

func IsUserExists() {

}

func UpdateContactList() {

}
