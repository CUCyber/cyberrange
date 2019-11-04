package server

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	ErrMatchFailed = errors.New("cyberrange: failed to match string")
)

func ParseUID(path string) (uint64, error) {
	var err error
	var uid uint64

	r := regexp.MustCompile(`/users/(?P<UID>\d+)$`)
	res := r.FindStringSubmatch(path)
	names := r.SubexpNames()

	if len(res) != len(names) {
		return 0, ErrMatchFailed
	}

	for idx, name := range names {
		if name == "UID" {
			uid, err = strconv.ParseUint(res[idx], 10, 64)
			if err != nil {
				return uid, err
			}
		}
	}

	return uid, nil
}
