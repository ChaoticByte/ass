package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)
// License: MIT

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/gliderlabs/ssh"
)

var sessionListMutex = &sync.Mutex{}
var sessionList = []*ssh.Session{}

func PubkeyAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	// verify public key for auth
	u := ctx.User()
	pubkey := config.Clients[u]
	if pubkey == "" { // username not in config.Clients
		log.Printf("Authentication failure: Unknown user %s", u)
		return false
	}
	k, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubkey))
	if err != nil {
		log.Printf("Authentication failure: Could not parse public key for %s", u)
	}
	if ssh.KeysEqual(k, key) {
		return true // bassd
	} else {
		log.Printf("Authentication failure: Invalid public key for user %s", u)
		return false
	}
}

func AddToSessionList(s *ssh.Session) int {
	// add new session to list
	var idx int = -1
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()
	lenClientList := len(sessionList)
	// try to reuse an existing nil position in the slice
	for i := range(lenClientList) {
		if sessionList[i] == nil {
			// we found one!
			sessionList[i] = s
			idx = i
			break
		}
	}
	if idx == -1 {
		// we have to append
		sessionList = append(sessionList, s)
		idx = lenClientList
	}
	return idx
}

func RemoveFromSessionList(idx int) {
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()
	// we have to set it to nil instead of constructing a new
	// slice, so that the other session handlers have still
	// the correct index for their own removal
	sessionList[idx] = nil // session associated with the pointer reference should be freed automatically
}

func BroadcastLine(line []byte) {
	if logFlag { log.Printf("%s", line) } // output message
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()
	for i := range(len(sessionList)) {
		session := sessionList[i]
		if session != nil {
			(*session).Write(line)
		}
	}
}

func HandleSession(session ssh.Session) {
	// get username
	user := session.Context().User()
	msgPrefix := []byte(fmt.Sprintf("%s: ", user))
	// Add client to list
	sessionListIdx := AddToSessionList(&session)
	defer func () { // cleanup tasks
		RemoveFromSessionList(sessionListIdx)
		log.Printf("User %s disconnected.\n", user)
		BroadcastLine([]byte(fmt.Sprintf("[disconnected] %s\n", user)))
	}()
	// helo
	log.Printf("User %s connected.\n", user)
	BroadcastLine([]byte(fmt.Sprintf("[connected] %s\n", user)))
	// IO
	linebuf := bufio.NewReader(session)
	for {
		line, err := linebuf.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("An error occurred! ", user, " ", err)
			}
			break // disconnect
		} else {
			BroadcastLine(slices.Concat(msgPrefix, line))
		}
	}
}

func RunServer() {
	log.Printf("Starting SSH server on %s:%v", config.Host, config.Port)
	log.Fatal(ssh.ListenAndServe(
		fmt.Sprintf("%s:%v", config.Host, config.Port),
		HandleSession,
		ssh.HostKeyFile(privateKeyFilepath),
		ssh.PublicKeyAuth(PubkeyAuth),
		ssh.NoPty()))
}
