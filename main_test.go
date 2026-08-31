package main

import (
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestLoadPosts(t *testing.T) {
	posts, legacy, err := loadPosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("no published posts rendered")
	}
	seen := map[string]bool{}
	for i, p := range posts {
		if p.Title == "" || p.Slug == "" || p.ISODate == "" || strings.TrimSpace(string(p.HTML)) == "" {
			t.Errorf("missing rendered post fields: %q", p.Slug)
		}
		if seen[p.Slug] {
			t.Errorf("duplicate slug: %s", p.Slug)
		}
		seen[p.Slug] = true
		if i > 0 && posts[i-1].ISODate < p.ISODate {
			t.Errorf("posts out of date order at %s", p.Slug)
		}
		if p.LegacyURL != "" && legacy[p.LegacyURL] != p {
			t.Errorf("missing legacy URL mapping for %s", p.Slug)
		}
	}
}

func TestFrontMatterDates(t *testing.T) {
	for _, value := range []string{"2026-08-31", "2026-08-31T12:30:00Z"} {
		var meta frontMatter
		if err := yaml.Unmarshal([]byte("title: Example\ndate: "+value+"\ncategories: [family, history]\ndraft: true\n"), &meta); err != nil {
			t.Fatal(err)
		}
		if meta.Date.Format("2006-01-02") != "2026-08-31" || meta.Title != "Example" || !meta.Draft || len(meta.Categories) != 2 {
			t.Errorf("front matter decoded incorrectly: %+v", meta)
		}
	}
	var meta frontMatter
	if err := yaml.Unmarshal([]byte("date: not-a-date\n"), &meta); err == nil {
		t.Fatal("invalid date accepted")
	}
}
