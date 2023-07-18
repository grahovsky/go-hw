package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	easyjson "github.com/mailru/easyjson"
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
	scaner := bufio.NewScanner(r)
	scaner.Split(bufio.ScanLines)

	i := 0
	for scaner.Scan() {
		if err = easyjson.Unmarshal(scaner.Bytes(), &result[i]); err != nil {
			return
		}
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, domain) {
			match := re.FindStringSubmatch(user.Email)
			if len(match) > 0 {
				result[strings.ToLower(match[1])]++
			}
		}
	}

	return result, nil
}
