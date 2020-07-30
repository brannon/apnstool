// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package apns

import "strconv"

type Expiration int
type Priority int

const (
	NoExpiration Expiration = -1
	NoPriority   Priority   = -1
)

func ParseExpiration(s string) (Expiration, error) {
	intValue, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return NoExpiration, err
	}
	return Expiration(int(intValue)), nil
}

func (e Expiration) String() string {
	return strconv.Itoa(int(e))
}

func ParsePriority(s string) (Priority, error) {
	intValue, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return NoPriority, err
	}
	return Priority(int(intValue)), nil
}

func (p Priority) String() string {
	return strconv.Itoa(int(p))
}
