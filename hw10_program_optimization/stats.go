package hw10programoptimization

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

//go:generate ffjson $GOFILE
type User struct {
	Email string
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
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	reg, err := regexp.Compile("\\." + domain + "$")
	if err != nil {
		return nil, err
	}

	result := make(DomainStat)

	for i := range u {
		user := u[i]
		if reg.MatchString(user.Email) {
			x := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[x]++
		}
	}

	return result, nil
}
