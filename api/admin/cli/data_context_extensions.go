package main

type PasswordManager interface {
	SetPassword(realmName string, userName string, password string) error
}
