package slowmgr

import (
	"errors"
	"fmt"
	stringPkg "github.com/sonemaro/slow/internal/util/string"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ParserPg struct{}

func NewParserPg() *ParserPg {
	return &ParserPg{}
}

const slowQueryRegexp = `\s*(?P<level>[A-Z0-9]+):\s+duration: (?P<duration>[0-9\.]+) ms\s+(?:(statement)|(execute \S+)): `

var ErrInvalidSlowQueryLog = errors.New("invalid slow query log entry")

// Parse tries to parse a new entry if it has all conditions of
// a slow query log. Here it is a sample slow query log:
// 2022-01-20 03:58:09.683 +0330 [5628] postgres@postgres LOG:  duration: 2002.358 ms  statement: select pg_sleep(2);
func (ppg *ParserPg) Parse(r io.Reader) (*Slow, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	regx, err := regexp.Compile(slowQueryRegexp)
	if err != nil {
		return nil, err
	}

	str := buf.String()
	strMap, slowParams := strMapFind(str, regx)
	if strMap == "" {
		return nil, ErrInvalidSlowQueryLog
	}

	duration, err := strconv.ParseFloat(slowParams["duration"], 64)
	if err != nil {
		return nil, errors.New("cannot convert duration to float. error: " + err.Error())
	}
	query, err := getStatementFromLog(str)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot find statement. error: %s", err.Error()))
	}

	op, err := extractOperation(query)
	if err != nil {
		return nil, errors.New("cannot extract operation. error: " + err.Error())
	}

	return &Slow{
		Query:     query,
		Operation: op,
		Duration:  duration,
	}, nil

}

func strMapFind(s string, r *regexp.Regexp) (string, map[string]string) {
	match := r.FindStringSubmatch(s)
	if match == nil {
		return "", nil
	}

	captures := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		if name != "" {
			// ignore unnamed matches
			captures[name] = match[i]
		}
	}
	return match[0], captures
}

func extractOperation(sql string) (string, error) {
	sql = strings.TrimSpace(strings.ToLower(sql))
	ops := []string{"insert", "update", "delete", "select"}
	for _, s := range ops {
		if strings.HasPrefix(sql, s) {
			return s, nil
		}
	}
	return "", errors.New("no valid operations found")
}

func getStatementFromLog(s string) (string, error) {
	if stm, empty := stringPkg.GetStringInBetween(s, "statement:", ";"); empty {
		return "", errors.New("cannot find a valid statement")
	} else {
		return strings.TrimSpace(stm), nil
	}
}
