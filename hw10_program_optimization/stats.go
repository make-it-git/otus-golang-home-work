package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//go:generate ffjson $GOFILE
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reg, err := regexp.Compile("@(?P<Domain>.+\\." + domain + ")$")
	if err != nil {
		return nil, fmt.Errorf("regexp error: %w", err)
	}
	reDomainIndex := reg.SubexpIndex("Domain")
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User
	for scanner.Scan() {
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}
		matches := reg.FindStringSubmatch(user.Email)
		if len(matches) > reDomainIndex {
			matchedDomain := strings.ToLower(matches[reDomainIndex])
			result[matchedDomain]++
		}
	}

	return result, nil
}
