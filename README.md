# jmainguy.com

Jonathan's personal site, delivered as one self-contained Go binary. The server,
templates, compiled TypeScript and Tailwind CSS, and Markdown posts are embedded
at build time. Blog media is hosted by Immich and tracked in
`immich-assets.json`.

## Add new posts
```sh
cp content/posts/example.md content/posts/my-first-post.md
```

## Preview changes
```sh
make dev
```


## Deploy changes
```sh
make build
./bin/jmainguy.com
```

The server listens on port `8080` by default. Set `PORT` to override it.
