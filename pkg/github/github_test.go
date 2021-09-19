package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const correctToken = "asdfASFA_245d1e5t2d1q6e5t3d2a55dffdhutf2"

func TestTokenCorrectness(t *testing.T) {
	data := []struct {
		name      string
		token     string
		isCorrect bool
	}{
		{name: "success", token: correctToken, isCorrect: true},
		{name: "missing token", token: "", isCorrect: false},
		{name: "short token", token: "sdfs", isCorrect: false},
		{name: "long token", token: "s4545422s5d5f5we55s5s5d5rt5f5s5s5dd225dfs", isCorrect: false},
		{name: "invalis token", token: "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$", isCorrect: false},
	}

	for _, d := range data {
		fn := func(t *testing.T) {
			_, err := NewClient("https://api.github.com", d.token)

			if !d.isCorrect && err == nil {
				t.Errorf("NewClient (%q), should error but did not", d.token)
			}

			if d.isCorrect && err != nil {
				t.Errorf("NewClient (%q), should't error but it did: %v", d.token, err)
			}
		}
		t.Run(d.name, fn)
	}
}

func TestContributorsList(t *testing.T) {
	f := func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Errorf("Method did not match: Got %q, want %q", r.Method, http.MethodGet)
			return
		}

		if got, want := r.URL.Path, "/repos/golang/go/contributors"; got != want {
			t.Errorf("Path did not match: Got %q want %q", got, want)
			return
		}

		if got, want := r.Header.Get("Authorization"), "token "+correctToken; got != want {
			t.Errorf("Token did not match: Got %q want %q", got, want)
			return
		}
		response := `[{"login":"anna", "contributions": 25}, {"login":"jacob", "contributions":12}]`
		if _, err := w.Write([]byte(response)); err != nil {
			t.Fatal(err)
		}
	}

	srv := httptest.NewServer(http.HandlerFunc(f))
	defer srv.Close()

	c, err := NewClient(srv.URL, correctToken)
	if err != nil {
		t.Fatal(err)
	}
	cons, err := c.ContributorList("golang/go")
	if err != nil {
		t.Fatalf("Client should not get error, but got %v", err)
	}

	want := []Contributor{
		{Login: "anna", Contributions: 25},
		{Login: "jacob", Contributions: 12},
	}

	if diff := cmp.Diff(cons, want); diff != "" {
		t.Errorf("contributors returned from client fif not match expected:\n%s", diff)
	}
}

func TestContributorsAPIFailure(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}
	srv := httptest.NewServer(http.HandlerFunc(fn))

	c, err := NewClient(srv.URL, correctToken)
	if err != nil {
		t.Errorf("Error during client creation")
	}

	_, err = c.ContributorList("golang/go")
	if err == nil {
		t.Fatal("Client should error but did not")
	}
}
