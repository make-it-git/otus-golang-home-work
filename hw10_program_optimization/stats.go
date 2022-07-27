package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"regexp"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reg, err := regexp.Compile("@(?P<Domain>.+\\." + domain + ")$")
	if err != nil {
		return nil, fmt.Errorf("regexp error: %w", err)
	}
	reDomainIndex := reg.SubexpIndex("Domain")
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		if email == "" {
			return nil, fmt.Errorf("'Email' field not found")
		}
		matches := reg.FindStringSubmatch(email)
		if len(matches) > reDomainIndex {
			matchedDomain := strings.ToLower(matches[reDomainIndex])
			result[matchedDomain]++
		}
	}

	return result, nil
}
