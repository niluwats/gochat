package service

import (
	"github.com/niluwats/gochat/pkg/dto"
	"github.com/niluwats/gochat/pkg/redisrepo"
)

func Register(u *dto.UserReq) *dto.Response {
	res := &dto.Response{
		Status: true,
	}

	status := redisrepo.IsUserExists(u.Username)
	if status {
		res.Status = false
		res.Message = "username already taken"
		return res
	}

	err := redisrepo.RegisterNewUser(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return res
	}

	return res
}

func Login(u *dto.UserReq) *dto.Response {
	res := &dto.Response{Status: true}

	err := redisrepo.IsUserAuthenticated(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return res
	}

	return res
}

func VerifyContact(username string) *dto.Response {
	res := &dto.Response{Status: true}
	status := redisrepo.IsUserExists(username)
	if !status {
		res.Status = false
		res.Message = "invalid username"
		return res
	}
	return res
}

func ChatHistory(username1, username2, fromTs, toTs string) *dto.Response {
	res := &dto.Response{}
	if !redisrepo.IsUserExists(username1) || !redisrepo.IsUserExists(username2) {
		res.Message = "incorrect username"
		return res
	}

	chats, err := redisrepo.FetchChatBetween(username1, username2, fromTs, toTs)
	if err != nil {
		res.Message = "unable to fetch history " + err.Error()
		return res
	}
	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func ContactList(username string) *dto.Response {
	res := &dto.Response{}

	if !redisrepo.IsUserExists(username) {
		res.Message = "incorrect username"
		return res
	}

	contactList, err := redisrepo.FetchContactList(username)
	if err != nil {
		res.Message = err.Error()
		return res
	}

	res.Data = contactList
	res.Status = true
	res.Total = len(contactList)
	return res
}
