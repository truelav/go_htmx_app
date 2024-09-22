package main

// import (
// 	"fmt"
// )

type authInfo struct {
	username string
	password string
}

// func (authI authInfo) getBasicAuth() string {
// 	return fmt.Sprintf("Authorisation Basic %s:%s", authInfo.username, authInfo.password)
// }

type employee interface {
	getName() string
	getSalary() int
}

type contractor struct {
	name         string
	hourlyPay    int
	hoursPerYear int
}

func (c contractor) getName() string {
	return c.name
}

func (c contractor) getSalary() int {
	return c.hourlyPay * c.hoursPerYear
}
