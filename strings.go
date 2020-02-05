package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"fmt"
	"sort"
	"strings"
)

type builder struct {
	list   []string
	str    string
	offset int
}

func (s *builder) add(entry string) {
	s.list = append(s.list, entry)
}

func (s *builder) write(w writer) error {
	var builder strings.Builder

	s.offset = w.offset()

	sort.Slice(s.list, func(i, j int) bool { return len(s.list[i]) > len(s.list[j]) })
	s.str = builder.String()

	for _, entry := range s.list {
		if !strings.Contains(s.str, entry) {
			builder.WriteString(entry)
			s.str = builder.String()
		}
	}

	_, err := w.Write([]byte(s.str))

	return err
}

func (s *builder) slice(entry string) string {
	if pos := strings.Index(s.str, entry); pos >= 0 && len(entry) > 0 {
		return fmt.Sprintf("/* %s */ str[%d:%d]", s.str[pos:pos+len(entry)], s.offset+pos, s.offset+pos+len(entry))
	}

	return `"` + entry + `"`
}
