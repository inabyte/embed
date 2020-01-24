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
	list    []string
	builder strings.Builder
	offset  int
}

func (s *builder) add(entry string) {
	s.list = append(s.list, entry)
}

func (s *builder) process() {
	sort.Slice(s.list, func(i, j int) bool { return len(s.list[i]) > len(s.list[j]) })

	for _, entry := range s.list {
		if !strings.Contains(s.builder.String(), entry) {
			s.builder.WriteString(entry)
		}
	}
}

func (s *builder) len() int64 {
	return int64(s.builder.Len())
}

func (s *builder) bytes() []byte {
	return []byte(s.builder.String())
}

func (s *builder) slice(entry string) string {
	if len(entry) == 0 {
		return `""`
	}

	pos := strings.Index(s.builder.String(), entry)

	if pos < 0 {
		return "\"" + entry + "\""
	}

	return fmt.Sprintf("/* %s */ str[%d:%d]", s.builder.String()[pos:pos+len(entry)], s.offset+pos, s.offset+pos+len(entry))
}
