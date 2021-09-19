package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/piotrstrzalka/contributors/pkg/github"
)

type mock struct{}

func (m *mock) ContributorList(s string) ([]github.Contributor, error) {
	switch {
	case s == "golang/go":
		return []github.Contributor{
			{Login: "tom", Contributions: 255},
			{Login: "johny", Contributions: 478},
			{Login: "alice", Contributions: 999},
		}, nil
	default:
		return nil, errors.New("could not reach API")
	}
}

func TestPrintContributors(t *testing.T) {
	data := []struct {
		name      string
		address   string
		want      string
		shouldErr bool
	}{
		{
			name:      "golang",
			address:   "golang/go",
			want:      "1. tom 255\n2. johny 478\n3. alice 999\n",
			shouldErr: false,
		},
		{
			name:      "error",
			address:   "balblabla",
			shouldErr: true,
		},
	}

	for _, d := range data {

		m := &mock{}
		fn := func(t *testing.T) {

			var buf bytes.Buffer
			err := printContributors(d.address, m, &buf, "")

			if d.shouldErr && err == nil {
				t.Fatal("There should be error but there wasn't")
			}

			if !d.shouldErr && err != nil {
				t.Fatalf("There shouldn't be error but there was: %v", err)
			}

			if d.want != buf.String() {
				t.Errorf("Got value is not the same as expected.\n")
				t.Logf("expected\n%s", d.want)
				t.Logf("printed\n%s", buf.String())
			}
		}

		t.Run(d.name, fn)
	}
}
