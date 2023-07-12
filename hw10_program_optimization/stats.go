package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email    string
	ID       int
	Name     string
	Username string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var re = regexp.MustCompile(`.*@(.*)`)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bufio.NewReader(r))
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

	for _, user := range u {
		if strings.HasSuffix(user.Email, domain) {
			match := re.FindStringSubmatch(user.Email)
			if len(match) > 0 {
				lowerDomain := strings.ToLower(match[1])
				result[lowerDomain]++
			}
		}
	}

	return result, nil
}
