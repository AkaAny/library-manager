package model

import (
	"github.com/go-session/session"
	"library-manager/logger"
)

type SessionInfo struct {
	Login    bool
	UserName string
	Working  bool
}

func IsLogin(obj interface{}) bool {
	info, ok := obj.(SessionInfo)
	if !ok {
		logger.Error.Printf("fail to convert %v -> SessionInfo", obj)
		return false
	}
	return info.Login
}

func GetFromStore(store session.Store) *SessionInfo {
	infoObj, exist := store.Get("info")
	if !exist {
		return nil
	}
	var info = infoObj.(SessionInfo)
	return &info
}
