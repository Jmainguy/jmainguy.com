---
title: "Rebuilding jmainguy.com as One Go Binary"
date: 2026-09-16T12:00:00-04:00
draft: false
image: "/images/social/rebuilding-jmainguy-com-as-one-go-binary.png"
categories:
  - go
  - web
  - kubernetes
  - gitops
---

A few years ago, the wise [Chris Short](https://www.linkedin.com/in/thechrisshort/) inspired me to convert my blog from WordPress to [Hugo](https://gohugo.io/about/introduction/).

[WordPress](https://wordpress.org/) had been nice in a lot of ways. It had a good editor, made uploading photos easy, and even allowed people to comment. What I did not enjoy was maintaining it. WordPress needed regular updates, plugins would break, themes did not always age well, and the site required a MySQL database and FTP access.

Moving to Hugo reduced that operational burden considerably. My posts became Markdown files stored in Git, which I really liked. The site no longer needed a database or a collection of plugins, and there was much less infrastructure to maintain.

However, I never found Hugo especially easy to use. Every time I wanted to publish something, I had to look up the documentation to remember how to create a post. Hugo also had themes, but I was never very happy with the ones I found. In my opinion, my site was still kind of ugly.

I kept the blog on Hugo for several years because it worked, but I was ready for a change.

<!-- SCREENSHOT: Add a side-by-side image of the old Hugo homepage and the new jmainguy.com homepage. -->

## Building the stack I wanted

Over the last year, I have been getting more interested in web design, with help from LLMs. I have settled on an architecture I really like: one small container running a single Go binary, with the TypeScript, [Tailwind CSS](https://tailwindcss.com/), HTML templates, Markdown posts, and other frontend assets embedded directly into it.

I like having a single container and binary for maintenance and security reasons. There are fewer moving pieces to update, deploy, and troubleshoot.

With that in mind, I designed the new [jmainguy.com](https://jmainguy.com/) using this stack. The blog posts are still Markdown files stored in Git, which was my favorite part of Hugo. I wrote the rest of the site in Go. The complete [source repository is public on GitHub](https://github.com/Jmainguy/jmainguy.com).

## What is inside the repository

The repository is small enough that I can understand the whole application:

```text
.
├── content/posts/       Markdown Logbook entries
├── static/              Images, icons, and the web manifest
├── web/
│   ├── templates/       Go HTML templates
│   ├── app.ts           Browser behavior
│   ├── styles.css       Tailwind source and site styles
│   └── dist/            Compiled CSS and JavaScript
├── main.go              HTTP server, routing, feeds, and post loading
├── projects.go          Project-page content
├── Makefile             Local build and preview commands
└── .goreleaser.yml      Binary and container release configuration
```

The frontend build compiles the TypeScript with [esbuild](https://esbuild.github.io/) and generates the CSS with Tailwind. Go's [`embed`](https://pkg.go.dev/embed) package then puts the finished frontend, templates, static files, and every Markdown post into the executable at build time:

```go
// Everything needed at runtime is compiled into the executable.
//
//go:embed content/posts web/dist web/templates static
var files embed.FS
```

When the server starts, it reads the embedded Markdown files, parses their YAML front matter, renders the Markdown with [Goldmark](https://github.com/yuin/goldmark), and prepares the posts for the Logbook. Go's standard `net/http` package handles the routes and serves the embedded CSS, JavaScript, templates, and images.

The request path is short:

```text
Browser request
    ↓
Go net/http router
    ↓
Embedded HTML template + rendered Markdown
    ↓
HTML response
```

Each request stays inside the Go application. The deployed artifact already contains the server, templates, frontend assets, and posts.

## What Battlestar Galactica taught me about containers

I originally chose this style of container because of something I learned from *Battlestar Galactica*: the only network the Cylons cannot hack is the one that does not exist.

I apply the same idea to containers. Software that is not present cannot be exploited. A scratch container starts with nothing except what you explicitly place inside it. There is no shell, package manager, or collection of utilities waiting to be maintained or abused. That limits the attack surface and makes the image easier to keep updated.

The current jmainguy.com image uses [Chainguard's `static` base image](https://images.chainguard.dev/directory/image/static/overview) through [ko](https://ko.build/). It follows the same minimal approach as `scratch` and provides a nonroot user and the small set of runtime files a static binary may need, such as CA certificates. I build the Go binary with `CGO_ENABLED=0`, which produces a statically linked binary.

The smaller image gives me less to maintain and less for an attacker to work with. The final image contains the site, its CA certificates, and the files required to run it.

## Making the site feel like mine

I generated the logo for the site to represent three important parts of my life: Christ in the cross, my family in the tree, and open source in the roots, which resemble Git branches.

I have always been partial to green, so the color scheme feels natural to me. More broadly, I wanted the site to feel personal. I wanted a place where I could publish writing, showcase my projects, and build fun things such as an interactive family tree.

Since I wrote the site, I can build whatever features I want into it.

I am also a big *Star Trek* fan. Each episode begins with something like, "Captain's log, stardate...," which I have always found fun. That inspired me to call my blog a **Logbook** and each post a **log entry**.

<!-- SCREENSHOT: Add the new homepage in light mode, showing the logo, recent logs, and selected projects. -->

## Responsive design and themes

I wanted the design to make good use of a large monitor while still working well on a phone. The layout starts with a narrow-screen design and uses [Tailwind's responsive breakpoints](https://tailwindcss.com/docs/responsive-design) to expand the grids, navigation, typography, and spacing as more room becomes available. I test the smallest layout at 360 by 616 pixels and the largest at 1920 by 1080.

The site also has light and dark themes. On the first visit it respects the operating system's preferred color scheme. The theme button lets the reader override that choice, and a small piece of TypeScript saves the preference in `localStorage`. Tailwind's [dark-mode variant](https://tailwindcss.com/docs/dark-mode) applies the alternate palette.

<!-- SCREENSHOT: Add mobile screenshots of the homepage in light and dark mode. -->

## Handling images

I already run an [Immich](https://immich.app/) server for my photos, so I decided to reuse it for the Logbook. I upload an image to Immich and reference its URL from the Markdown post.

That lets me keep the writing itself in portable Markdown without adding a second image-management system to the site. The site also gives article images a lightbox, so a reader can open them at a larger size and move between all the images in a post.

## The Mainguy family tree

For the family tree, I am building on the hard work of Bill Mainguy, a distant relative of mine. Bill created [mainguy.ca](http://mainguy.ca/) and compiled an impressive amount of information about the Mainguy family.

I love finding passion projects like his. His work has been invaluable to me because I am very interested in my family history, and much of this information would otherwise be difficult to find. Bill made it much easier to understand where one side of my family came from.

He has created quite a legacy. I hope that one day I can make something as useful and meaningful as what Bill made.

Using the information he compiled, I created my own visualization of the [Mainguy family tree](https://jmainguy.com/family/). I have included my wife and children, but I still need to get permission from my siblings and cousins before adding their immediate families. That will be one of my next steps.

<!-- SCREENSHOT: Add the family tree focused on Holden and Hazel. -->

## A home for my projects

The site also gives me a place to showcase the things I build. The [projects page](https://jmainguy.com/projects/) collects some of them in one easy-to-find place, including [Shoreline](https://jmainguy.com/projects/shoreline/), [Verbose Resume](https://jmainguy.com/projects/verbose-resume/), and [soh.re](https://jmainguy.com/projects/soh-re/).

I love working on projects like these, and I plan to keep adding to that page as I build new things and revisit older ones.

<!-- SCREENSHOT: Add the projects page or a three-project montage. -->

## Building SEO into the site

Owning the site means I also own the parts that a framework or plugin used to generate for me. Each page includes a canonical URL, description, Open Graph metadata, Twitter card metadata, and structured data from [Schema.org](https://schema.org/). The server produces a [`robots.txt`](https://jmainguy.com/robots.txt), [`sitemap.xml`](https://jmainguy.com/sitemap.xml), and [RSS feed](https://jmainguy.com/feed.xml) directly from the same embedded posts and project data.

The old WordPress and Hugo URLs redirect to their matching Logbook entries, which preserves old links and bookmarks. Internal links connect related writing, such as my posts about [building a good CI system](/logbook/a-good-ci-system/), [GitOps](/logbook/gitops/), and [running a multi-architecture k3s cluster](/logbook/multi-arch-hybrid-cloud/).

I also added an [`llms.txt`](https://jmainguy.com/llms.txt) endpoint. It gives language models a plain-text overview of the main pages, projects, and recent Logbook entries, with links back to the original pages. [`llms.txt`](https://llmstxt.org/) is still a proposal, but it is inexpensive to generate from data the application already has.

## Writing with an LLM and keeping my voice

I am still refining my writing process for the Logbook.

For this post, I wrote the original text myself and pasted it into Codex. I asked GPT-5.6 Sol to fix my typos, help organize the post, prompt me to fill in missing sections, and point out places that were difficult to read or needed more detail.

I am cautious about having an LLM write the post for me because that can remove my voice and the way I naturally speak. For now, I am trying to limit its role to editing, structure, and feedback.

I do not consider myself a great storyteller. I love listening to stories, but when I tell one, I tend to get straight to the point. In doing that, I sometimes skip the details that make a story enjoyable. I am hoping an LLM can help identify those missing details while leaving the actual story and words to me.

The intended reader is often a future version of myself. I want that person to recognize the voice as mine.

## From Markdown to production

Once I have a post where I like it, I run the site locally with `make dev` and open it in my browser. I read through the post, check the responsive layout and themes, and make sure the images look right.

When I am happy with it, I push the change to Git. That starts a shared [GitHub Actions](https://github.com/features/actions) release workflow I use across my Go projects:

1. [Release Please](https://github.com/googleapis/release-please) creates or updates a release pull request from my conventional commits.
2. Merging that pull request creates the GitHub release.
3. [GoReleaser](https://goreleaser.com/) builds the Linux binaries for AMD64 and ARM64.
4. ko builds the container images and publishes them to my [zot OCI registry](https://zotregistry.dev/).
5. The release workflow signs the published images with [Notation](https://notaryproject.dev/docs/user-guides/installation/cli/).
6. [Argo CD](https://argo-cd.readthedocs.io/en/stable/) sees the desired version in Git and deploys it to my [k3s](https://docs.k3s.io/) Kubernetes cluster.

```text
Markdown + code
      ↓
Git push
      ↓
Release Please
      ↓
GoReleaser + ko
      ↓
Signed image in zot
      ↓
Argo CD
      ↓
k3s cluster
```

I wrote more about the deployment side in [GitOps: Transforming Deployment with Automation and Version Control](/logbook/gitops/) and about the cluster itself in [Multi-Arch Hybrid Cloud](/logbook/multi-arch-hybrid-cloud/).

## The tradeoff

I am responsible for routing, Markdown rendering, templates, feeds, metadata, redirects, responsive behavior, accessibility, and security headers. When I want a new feature, I have to understand it and add it myself.

I am happy with that trade. The application is small, its dependencies are limited, and the whole runtime fits into one binary and one minimal container. I can trace a request from the router to the template. I can add unusual pages, such as the interactive family tree, using whatever data and layout make sense for the project.

Publishing begins with a Markdown file in Git. Local previewing is one command. The same tests and release process I use for my other Go projects handle the site. I enjoy working on it, which makes me more likely to write and publish.

## A site that feels like me

I really love the new design. It feels more like me than either the Hugo or WordPress versions did.

I like the simplicity of running the site as a single Go binary. I like that the posts remain Markdown files stored in Git. I like being able to build pages that reflect my faith, family, work, and interests. I also like being able to work directly on the site with LLMs and shape it into whatever I want it to be.

For the first time, I feel like I built my own place on the web.
