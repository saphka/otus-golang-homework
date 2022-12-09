package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	//ID       int
	//Name     string
	//Username string
	Email string
	//Phone    string
	//Password string
	//Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	buf := bufio.NewReader(r)
	domWithDot := "." + domain
	var user User

	for shouldExit := false; !shouldExit; {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			shouldExit = true
		} else if err != nil {
			return nil, err
		}

		if err = easyjson.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domWithDot) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
