package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	decoder := json.NewDecoder(r)
	i := 0
	for decoder.More() {
		var user User
		if err = decoder.Decode(&user); err != nil {
			return result, err
		}
		result[i] = user
		i++
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	re, err := regexp.Compile(`.*@(.*.` + domain + `)`)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		match := re.FindStringSubmatch(user.Email)
		if len(match) > 0 {
			result[strings.ToLower((match[1]))]++
		}
	}

	return result, nil
}
