// Copyright 2019 George Aristy
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commits

import (
	"fmt"
	"io"
	"strings"

	"github.com/llorllale/go-gitlint/internal/repo"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Commits returns commits.
// @todo #4 Figure out how to disable the golint check that
//  forces us to write redundant comments of the form
//  'comment on exported type Commits should be of the form
//  "Commits ..." (with optional leading article)' and rewrite
//  all comments.
type Commits func() []*Commit

// Commit holds data for a single git commit.
type Commit struct {
	hash    string
	message string
}

// Hash is the commit's ID.
func (c *Commit) Hash() string {
	return c.hash
}

// Subject is the commit message's subject line.
func (c *Commit) Subject() string {
	return strings.Split(c.message, "\n\n")[0]
}

// Body is the commit message's body.
func (c *Commit) Body() string {
	body := ""
	parts := strings.Split(c.message, "\n\n")
	if len(parts) > 1 {
		body = strings.Join(parts[1:], "")
	}
	return body
}

// In returns the commits in the repo.
// @todo #4 These err checks are extremely annoying. Figure out
//  how to handle them elegantly and reduce the cyclo complexity
//  of this function (currently at 4).
func In(repository repo.Repo) Commits {
	return func() []*Commit {
		r := repository()
		ref, err := r.Head()
		if err != nil {
			panic(err)
		}
		iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			panic(err)
		}
		commits := make([]*Commit, 0)
		_ = iter.ForEach(func(c *object.Commit) error { //nolint[errcheck]
			commits = append(
				commits,
				&Commit{
					hash:    c.Hash.String(),
					message: c.Message,
				},
			)
			return nil
		})
		return commits
	}
}

// Printed prints the commits to the file.
func Printed(commits Commits, writer io.Writer, sep string) Commits {
	return func() []*Commit {
		input := commits()
		for _, c := range input {
			_, err := writer.Write(
				[]byte(fmt.Sprintf("%s%s", &pretty{c}, sep)),
			)
			if err != nil {
				panic(err)
			}
		}
		return input
	}
}

// a Stringer for pretty-printing the commit.
type pretty struct {
	*Commit
}

func (p *pretty) String() string {
	return fmt.Sprintf(
		"hash: %s subject=%s body=%s",
		p.Commit.Hash(), p.Commit.Subject(), p.Commit.Body(),
	)
}