package main

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"go.yaml.in/yaml/v3"
)

// Everything needed at runtime is compiled into the executable.
//
//go:embed content/posts web/dist web/templates static
var files embed.FS

type frontMatter struct {
	Title      string   `yaml:"title"`
	Date       date     `yaml:"date"`
	Draft      bool     `yaml:"draft"`
	URL        string   `yaml:"url"`
	Categories []string `yaml:"categories"`
}

type date struct{ time.Time }

func (d *date) UnmarshalYAML(node *yaml.Node) error {
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if value, err := time.Parse(layout, node.Value); err == nil {
			d.Time = value
			return nil
		}
	}
	return fmt.Errorf("unsupported date %q", node.Value)
}

type post struct {
	Title, Slug, Excerpt, Search, ISODate, ShortDate, LongDate, RSSDate string
	LegacyURL                                                           string
	Categories                                                          []string
	HTML                                                                template.HTML
}

type postGroup struct {
	Year  string
	Posts []*post
}

type pageData struct {
	MetaTitle, Description, Canonical string
	Posts                             []*post
	PostGroups                        []postGroup
	Post                              *post
	Projects                          []*project
	Project                           *project
}

var youtube = regexp.MustCompile(`\{\{<\s*youtube\s+([^ >]+)\s*>\}\}`)
var tags = regexp.MustCompile(`<[^>]+>`)
var spaces = regexp.MustCompile(`\s+`)

func main() {
	posts, legacy, err := loadPosts()
	if err != nil {
		log.Fatal(err)
	}
	postGroups := groupPostsByYear(posts)
	templates := map[string]*template.Template{}
	for _, name := range []string{"home.html", "about.html", "family.html", "projects.html", "project.html", "posts.html", "post.html", "not_found.html"} {
		templates[name] = template.Must(template.ParseFS(files, "web/templates/base.html", "web/templates/"+name))
	}
	mux := http.NewServeMux()

	serveDir(mux, "/assets/", "web/dist")
	serveDir(mux, "/images/", "static/images")
	mux.HandleFunc("/robots.txt", serveRobots)
	mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, _ *http.Request) { serveSitemap(w, posts, projects) })
	mux.HandleFunc("/feed.xml", func(w http.ResponseWriter, _ *http.Request) { serveFeed(w, posts) })
	mux.HandleFunc("/llms.txt", func(w http.ResponseWriter, _ *http.Request) { serveLLMs(w, posts, projects) })
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		data, err := files.ReadFile("static/images/favicon.ico")
		if err != nil {
			http.NotFound(w, nil)
			return
		}
		w.Header().Set("Content-Type", "image/x-icon")
		_, _ = w.Write(data)
	})
	mux.HandleFunc("/site.webmanifest", func(w http.ResponseWriter, _ *http.Request) {
		data, err := files.ReadFile("static/site.webmanifest")
		if err != nil {
			http.NotFound(w, nil)
			return
		}
		w.Header().Set("Content-Type", "application/manifest+json")
		_, _ = w.Write(data)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			render(w, templates["not_found.html"], http.StatusNotFound, pageData{MetaTitle: "Page not found · Jonathan Mainguy", Description: "Page not found", Canonical: r.URL.Path})
			return
		}
		featured := posts
		if len(featured) > 7 {
			featured = featured[:7]
		}
		render(w, templates["home.html"], http.StatusOK, pageData{MetaTitle: "Jonathan Mainguy · Software Engineer", Description: "Jonathan Mainguy writes about faith, family, Linux, Go, Kubernetes, open source software, and leadership.", Canonical: "/", Posts: featured})
	})
	mux.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/about/" {
			render(w, templates["not_found.html"], http.StatusNotFound, pageData{MetaTitle: "Page not found · Jonathan Mainguy", Description: "Page not found", Canonical: r.URL.Path})
			return
		}
		render(w, templates["about.html"], http.StatusOK, pageData{MetaTitle: "About · Jonathan Mainguy", Description: "Jonathan Mainguy on faith in Jesus Christ, family, open source software, and leadership.", Canonical: "/about/"})
	})
	mux.HandleFunc("/family/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/family/" {
			render(w, templates["not_found.html"], http.StatusNotFound, pageData{MetaTitle: "Page not found · Jonathan Mainguy", Description: "Page not found", Canonical: r.URL.Path})
			return
		}
		render(w, templates["family.html"], http.StatusOK, pageData{MetaTitle: "Family · Jonathan Mainguy", Description: "An interactive Mainguy family tree spanning the eight-generation chart documented by Mainguy.ca and Jonathan Mainguy's family.", Canonical: "/family/"})
	})
	mux.HandleFunc("/logbook/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/logbook/"), "/")
		if slug == "" {
			render(w, templates["posts.html"], http.StatusOK, pageData{MetaTitle: "Logbook · Jonathan Mainguy", Description: "Jonathan Mainguy's logbook covering Linux, Go, Kubernetes, infrastructure, leadership, and life.", Canonical: "/logbook/", Posts: posts, PostGroups: postGroups})
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}
		for _, p := range posts {
			if p.Slug == slug {
				render(w, templates["post.html"], http.StatusOK, pageData{MetaTitle: p.Title + " · Jonathan Mainguy", Description: p.Excerpt, Canonical: "/logbook/" + p.Slug + "/", Post: p})
				return
			}
		}
		render(w, templates["not_found.html"], http.StatusNotFound, pageData{MetaTitle: "Page not found · Jonathan Mainguy", Description: "Page not found", Canonical: r.URL.Path})
	})
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		target := "/logbook/" + strings.TrimPrefix(r.URL.Path, "/posts/")
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})
	mux.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/projects/"), "/")
		if slug == "" {
			render(w, templates["projects.html"], http.StatusOK, pageData{MetaTitle: "Projects · Jonathan Mainguy", Description: "Web projects, games, tools, and experiments built by Jonathan Mainguy.", Canonical: "/projects/", Projects: projects})
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}
		if project := projectBySlug(slug); project != nil {
			render(w, templates["project.html"], http.StatusOK, pageData{MetaTitle: project.Name + " · Projects · Jonathan Mainguy", Description: project.Summary, Canonical: "/projects/" + project.Slug + "/", Project: project})
			return
		}
		render(w, templates["not_found.html"], http.StatusNotFound, pageData{MetaTitle: "Page not found · Jonathan Mainguy", Description: "Page not found", Canonical: r.URL.Path})
	})
	for legacyPath, p := range legacy {
		target := "/logbook/" + p.Slug + "/"
		mux.Handle(legacyPath, http.RedirectHandler(target, http.StatusMovedPermanently))
	}

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("jmainguy.com listening on http://localhost%s (%d posts)", addr, len(posts))
	log.Fatal(http.ListenAndServe(addr, securityHeaders(mux)))
}

func groupPostsByYear(posts []*post) []postGroup {
	groups := make([]postGroup, 0)
	for _, p := range posts {
		year := p.ISODate[:4]
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, postGroup{Year: year})
		}
		groups[len(groups)-1].Posts = append(groups[len(groups)-1].Posts, p)
	}
	return groups
}

func loadPosts() ([]*post, map[string]*post, error) {
	entries, err := fs.ReadDir(files, "content/posts")
	if err != nil {
		return nil, nil, err
	}
	mdParser := parser.New(parser.WithExtensions(extension.GFMParser), parser.WithAutoHeadingID())
	mdRenderer := html.New(html.WithExtensions(extension.GFMHTMLRenderer), html.WithUnsafe())
	var posts []*post
	legacy := map[string]*post{}
	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".md" {
			continue
		}
		raw, err := files.ReadFile("content/posts/" + entry.Name())
		if err != nil {
			return nil, nil, err
		}
		parts := bytes.SplitN(raw, []byte("---"), 3)
		if len(parts) != 3 {
			return nil, nil, fmt.Errorf("%s: missing YAML front matter", entry.Name())
		}
		var meta frontMatter
		if err := yaml.Unmarshal(parts[1], &meta); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		if meta.Draft {
			continue
		}
		body := youtube.ReplaceAll(parts[2], []byte(`<iframe src="https://www.youtube-nocookie.com/embed/$1" title="YouTube video" loading="lazy" allowfullscreen></iframe>`))
		var rendered bytes.Buffer
		if err := mdRenderer.Render(&rendered, body, mdParser.Parse(body)); err != nil {
			return nil, nil, err
		}
		plain := strings.TrimSpace(spaces.ReplaceAllString(tags.ReplaceAllString(rendered.String(), " "), " "))
		runes := []rune(plain)
		if len(runes) > 190 {
			plain = string(runes[:190]) + "…"
		}
		slug := strings.TrimSuffix(entry.Name(), ".md")
		if len(slug) > 11 && slug[4] == '-' && slug[7] == '-' && slug[10] == '-' {
			slug = slug[11:]
		}
		p := &post{Title: meta.Title, Slug: slug, Excerpt: plain, HTML: template.HTML(rendered.String()), LegacyURL: meta.URL, Categories: meta.Categories, ISODate: meta.Date.Format("2006-01-02"), ShortDate: strings.ToUpper(meta.Date.Format("02 Jan 2006")), LongDate: meta.Date.Format("January 2, 2006"), RSSDate: meta.Date.Format(time.RFC1123Z)}
		p.Search = strings.ToLower(strings.Join([]string{p.Title, p.Excerpt, p.ISODate, strings.Join(p.Categories, " ")}, " "))
		posts = append(posts, p)
		if meta.URL != "" {
			legacy[meta.URL] = p
		}
	}
	sort.Slice(posts, func(i, j int) bool { return posts[i].ISODate > posts[j].ISODate })
	return posts, legacy, nil
}

func serveRobots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, "User-agent: *\nAllow: /\n\nSitemap: https://jmainguy.com/sitemap.xml\n")
}

func serveSitemap(w http.ResponseWriter, posts []*post, projects []*project) {
	type sitemapURL struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod,omitempty"`
	}
	type sitemap struct {
		XMLName xml.Name     `xml:"urlset"`
		Xmlns   string       `xml:"xmlns,attr"`
		URLs    []sitemapURL `xml:"url"`
	}
	urls := []sitemapURL{{Loc: "https://jmainguy.com/"}, {Loc: "https://jmainguy.com/about/"}, {Loc: "https://jmainguy.com/family/"}, {Loc: "https://jmainguy.com/projects/"}, {Loc: "https://jmainguy.com/logbook/"}}
	for _, p := range posts {
		urls = append(urls, sitemapURL{Loc: "https://jmainguy.com/logbook/" + p.Slug + "/", LastMod: p.ISODate})
	}
	for _, project := range projects {
		urls = append(urls, sitemapURL{Loc: "https://jmainguy.com/projects/" + project.Slug + "/"})
	}
	data, err := xml.MarshalIndent(sitemap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls}, "", "  ")
	if err != nil {
		http.Error(w, "sitemap error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write(append([]byte(xml.Header), data...))
}

func serveFeed(w http.ResponseWriter, posts []*post) {
	type rssItem struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		GUID        string `xml:"guid"`
		PubDate     string `xml:"pubDate"`
		Description string `xml:"description"`
	}
	type channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		Language      string    `xml:"language"`
		LastBuildDate string    `xml:"lastBuildDate"`
		Items         []rssItem `xml:"item"`
	}
	type rss struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel channel  `xml:"channel"`
	}
	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		link := "https://jmainguy.com/logbook/" + p.Slug + "/"
		items = append(items, rssItem{Title: p.Title, Link: link, GUID: link, PubDate: p.RSSDate, Description: p.Excerpt})
	}
	lastBuild := ""
	if len(posts) > 0 {
		lastBuild = posts[0].RSSDate
	}
	data, err := xml.MarshalIndent(rss{Version: "2.0", Channel: channel{Title: "Jonathan Mainguy's Logbook", Link: "https://jmainguy.com/", Description: "Logs about Linux, Go, Kubernetes, homelabs, software operations, leadership, and life.", Language: "en-us", LastBuildDate: lastBuild, Items: items}}, "", "  ")
	if err != nil {
		http.Error(w, "feed error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	_, _ = w.Write(append([]byte(xml.Header), data...))
}

func serveLLMs(w http.ResponseWriter, posts []*post, projects []*project) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintln(w, "# Jonathan Mainguy")
	_, _ = fmt.Fprintln(w, "\n> Personal website of Jonathan Mainguy, a software engineer in Garner, North Carolina. The site covers Linux, Go, Kubernetes, homelabs, infrastructure, and personal software projects.")
	_, _ = fmt.Fprintln(w, "\n## Main pages\n\n- [Home](https://jmainguy.com/): Biography and recent log entries.\n- [About](https://jmainguy.com/about/): Jonathan's faith, family, open source values, and approach to leadership.\n- [Family](https://jmainguy.com/family/): Jonathan's direct Mainguy family branch.\n- [Logbook](https://jmainguy.com/logbook/): Complete article archive.\n- [Projects](https://jmainguy.com/projects/): Websites, games, and tools built and operated by Jonathan.\n- [RSS feed](https://jmainguy.com/feed.xml): Machine-readable log feed.")
	_, _ = fmt.Fprintln(w, "\n## Projects")
	for _, project := range projects {
		_, _ = fmt.Fprintf(w, "\n- [%s](https://jmainguy.com/projects/%s/): %s", project.Name, project.Slug, project.Summary)
	}
	_, _ = fmt.Fprintln(w, "\n## Recent articles")
	limit := len(posts)
	if limit > 25 {
		limit = 25
	}
	for _, p := range posts[:limit] {
		_, _ = fmt.Fprintf(w, "\n- [%s](https://jmainguy.com/logbook/%s/): %s", p.Title, p.Slug, p.Excerpt)
	}
	_, _ = fmt.Fprintln(w, "\n\n## Attribution\n\nWhen citing this site, attribute the work to Jonathan Mainguy and link to the specific article or project page.")
}

func serveDir(mux *http.ServeMux, route, dir string) {
	sub, err := fs.Sub(files, dir)
	if err != nil {
		return
	}
	mux.Handle(route, http.StripPrefix(strings.TrimSuffix(route, "/"), http.FileServer(http.FS(sub))))
}

func render(w http.ResponseWriter, templates *template.Template, status int, data pageData) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		log.Printf("render: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self'; img-src 'self' https: data:; frame-src https://www.youtube-nocookie.com; font-src 'self'; connect-src 'self'")
		next.ServeHTTP(w, r)
	})
}
