package main

type projectLink struct {
	Name, URL string
}

type project struct {
	Name, Slug, Category, Summary, Image, ImageAlt, URL, Domain string
	Story                                                       []string
	Highlights                                                  []string
	Stack                                                       []string
	Repositories                                                []projectLink
}

var projects = orderProjects([]*project{
	{
		Name: "Shoreline", Slug: "shoreline", Category: "Generative art", URL: "https://shoreline.soh.re", Domain: "shoreline.soh.re",
		Image: "/images/projects/shoreline.webp", ImageAlt: "Moonlit procedural beach scene from Shoreline",
		Summary: "A procedural beach screensaver with tides, stars, moon phases, weather, and wandering seagulls.",
		Story: []string{
			"Shoreline is a quiet place to leave open on a second screen. The beach changes continuously as waves reach the sand, daylight gives way to night, stars move across the sky, and seagulls arrive and leave.",
			"The scene is generated in the browser without image assets. Astronomical data, canvas drawing, and simulation code create the sky, moon, ocean, sand, foam, and birds each time it runs.",
		},
		Highlights: []string{"A complete day passes every 15 real-time minutes", "Hipparcos star data and calculated local sidereal time", "Moon phases, earthshine, and embedded lunar crater data", "Procedural waves, reflections, foam, sand, and seagull behavior", "Reduced-motion support and optional ambient audio"},
		Stack:      []string{"Go", "JavaScript", "HTML Canvas", "Embedded data", "One binary"},
	},
	{
		Name: "BaboViolent", Slug: "baboviolent", Category: "Game and community", URL: "https://baboviolent.net", Domain: "baboviolent.net",
		Image: "/images/projects/baboviolent.webp", ImageAlt: "BaboViolent multiplayer arena",
		Summary: "The community home and browser-playable version of the cult multiplayer arena shooter.",
		Story: []string{
			"BaboViolent.net gives Babo Violent 2 players a home and keeps the game available to play. The project brings the original maps, models, textures, sounds, and fast arena combat into a public WebGL browser client.",
			"Bringing the game back took much more than wrapping the old executable. We treated the original C++ source as the specification and ported its behavior into plain JavaScript, WebGL2, and WebAudio. The browser client decodes the original TGA textures, BVM maps, and DKO models without modifying the game assets. We reproduced the movement, collision quirks, rolling babos, weapons, projectiles, particles, audio, HUD, game modes, and network protocol closely enough to preserve the way Babo feels.",
			"Multiplayer required a second major port. We rebuilt the authoritative C++ server rules in Rust and replaced the old BaboNet transport with binary WebSocket connections that browsers can use directly. The relay replaces the legacy master server with a modern server directory, authentication, and lobby backed by Redis and PostgreSQL. The website, containers, automated releases, Kubernetes deployment, public bot API, and Flagship CTF bot complete the system around the game.",
		},
		Highlights: []string{"Playable WebGL2 browser client", "Original unmodified game assets", "All 81 shipped maps and 26 shipped models", "Authoritative Rust multiplayer server", "Binary WebSocket game protocol", "Capture the Flag, King of the Hill, and bot support", "Server directory, accounts, lobby, and release automation"},
		Stack:      []string{"JavaScript", "WebGL2", "WebAudio", "Rust", "Go", "WebSockets", "Redis", "PostgreSQL", "Kubernetes"},
		Repositories: []projectLink{
			{Name: "Browser client", URL: "https://github.com/BaboViolent/baboviolent2-jsclient"},
			{Name: "Original Babo Violent 2 source", URL: "https://github.com/Daivuk/BaboViolent2"},
		},
	},
	{
		Name: "soh.re", Slug: "soh-re", Category: "Interactive terminal", URL: "https://soh.re", Domain: "soh.re",
		Image: "/images/projects/soh.gif", ImageAlt: "soh.re web terminal in action",
		Summary: "An interactive Linux terminal on the web, stocked with a few of my favorite command-line tools.",
		Story: []string{
			"soh.re opens directly into a real Linux shell in the browser. It is a small, hands-on introduction to my open source profile and a place to explore a few command-line programs without installing anything.",
			"Each session runs inside a Fedora environment exposed through ttyd. The welcome message points visitors toward the files in my profile and commands including cmatrix, asciiquarium, bible, and bak.",
			"The current version runs on Kubernetes. Istio routes regular HTTP traffic to the static upgrade service and sends WebSocket connections to isolation-proxy, a Go-based Kubernetes operator I wrote. The operator maintains a warm pool of ready pods, gives each connection its own fresh pod, blocks outbound internet access with a NetworkPolicy, and deletes the pod when the session closes.",
		},
		Highlights: []string{"Writable browser terminal", "One isolated pod for every connection", "Warm pool for fast session startup", "Automatic cleanup when sessions close", "Restricted outbound network access", "Built-in terminal tools and personal command-line projects"},
		Stack:      []string{"Fedora Linux", "ttyd", "Bash", "Go", "Containers", "Kubernetes", "Istio", "NetworkPolicy"},
		Repositories: []projectLink{
			{Name: "Terminal environment and container image", URL: "https://github.com/Jmainguy/soh.re"},
			{Name: "Kubernetes isolation proxy", URL: "https://github.com/Jmainguy/isolation-proxy"},
			{Name: "Helm and Istio deployment", URL: "https://github.com/jmainguy/helm-charts/tree/main/sohre"},
		},
	},
	{
		Name: "Bible Reader", Slug: "bible-reader", Category: "Scripture reader", URL: "https://bible.soh.re", Domain: "bible.soh.re",
		Image: "/images/projects/bible.webp", ImageAlt: "Bible Reader interface",
		Summary: "A fast Bible reader with direct verse links, multiple translations, and keyboard and touch navigation.",
		Story: []string{
			"Bible Reader makes Scripture easy to open, navigate, and share. A link can point directly to a book, chapter, verse, or verse range, which lets the rest of my sites cite the Bible from a reader I operate.",
			"The reader uses SWORD and CrossWire zText source material and serves several translations through a Go API. Chapter navigation works through visible controls, keyboard arrows, and mobile gestures.",
		},
		Highlights: []string{"Old and New Testament navigation", "Direct links to verses and verse ranges", "Multiple Bible translations", "Previous, next, selector, keyboard, and touch navigation", "Responsive reading interface"},
		Stack:      []string{"Go", "JavaScript", "Tailwind CSS", "SWORD and CrossWire data", "Embedded frontend"},
	},
	{
		Name: "Verbose Resume", Slug: "verbose-resume", Category: "Career tools", URL: "https://verboseresume.com", Domain: "verboseresume.com",
		Image: "/images/projects/verbose-resume.webp", ImageAlt: "Verbose Resume career profile and resume builder",
		Summary: "Keep one detailed career record, tailor it for every role, and track each application from one place.",
		Story: []string{
			"Verbose Resume keeps the complete version of a career in one structured record. Instead of repeatedly cutting down a master document by hand, users can select the experience that matters for a role and build a focused resume from it.",
			"The application also supports job imports, application tracking, public career profiles, tailored resume versions, and polished HTML output that prints cleanly to PDF.",
		},
		Highlights: []string{"One detailed career record", "Tailored resumes for individual roles", "Job imports and application tracking", "Public career profiles", "Print-ready PDF output"},
		Stack:      []string{"Go", "TypeScript", "Tailwind CSS", "PostgreSQL", "Embedded frontend"},
	},
	{
		Name: "NineVoltNine", Slug: "ninevoltnine", Category: "Music software", URL: "https://ninevoltnine.com", Domain: "ninevoltnine.com",
		Image: "/images/projects/ninevolt.webp", ImageAlt: "Ninevolt terminal guitar pedalboard running a live effects chain",
		Summary: "A Linux-first guitar pedalboard with real-time effects, keyboard controls, MIDI support, and shareable YAML chains.",
		Story: []string{
			"Ninevolt turns a Linux computer into a complete guitar pedalboard. Effects run in real time over JACK or PipeWire's JACK bridge, while a terminal interface provides readable knobs, bypass controls, presets, and live signal status.",
			"Pedals and complete chains can be saved as plain YAML, versioned, shared, and loaded without an account or internet connection. Keyboard and MIDI controls make the same setup useful at a desk or on stage.",
		},
		Highlights: []string{"Real-time guitar effects", "Terminal pedal controls", "JACK and PipeWire audio", "Keyboard and MIDI foot-controller support", "Shareable YAML pedalboards", "Offline operation"},
		Stack:      []string{"Rust", "Go", "TypeScript", "Tailwind CSS", "JACK", "PipeWire", "MIDI"},
	},
	{
		Name: "Heyyyeyaaeyaaaeyaeyaa", Slug: "heyyyeyaaeyaaaeyaeyaa", Category: "Internet history", URL: "https://hey.soh.re", Domain: "hey.soh.re",
		Image: "/images/projects/heman.webp", ImageAlt: "He-Man singing What's Up",
		Summary: "A permanent home for the He-Man cover of “What’s Up?” and a small piece of internet history.",
		Story: []string{
			"Some websites can have one job. hey.soh.re preserves a beloved piece of early internet culture: the animated He-Man performance of 4 Non Blondes' “What’s Up?”",
			"The page loads the video directly and keeps the experience focused on the thing people came to see. It is deliberately small, quick, and easy to remember.",
		},
		Highlights: []string{"One memorable address", "Direct video playback", "Minimal interface", "Self-contained deployment"},
		Stack:      []string{"HTML", "Video", "Container image", "Static site"},
	},
})

func orderProjects(catalog []*project) []*project {
	desired := []string{"baboviolent", "verbose-resume", "ninevoltnine", "soh-re", "bible-reader", "shoreline"}
	ordered := make([]*project, 0, len(catalog))
	for _, slug := range desired {
		for _, candidate := range catalog {
			if candidate.Slug == slug {
				ordered = append(ordered, candidate)
				break
			}
		}
	}
	return ordered
}

func projectBySlug(slug string) *project {
	for _, candidate := range projects {
		if candidate.Slug == slug {
			return candidate
		}
	}
	return nil
}
