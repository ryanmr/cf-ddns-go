# cf-ddns-go

> 👏 This a conversion of the existing [cf-ddns](https://github.com/ryanmr/cf-ddns), now written in [Go](https://go.dev/).

> 🟠 This is a work in progress.

## Usage

```sh
go run ./cmd/cli
```

```sh
docker compose up --build
```

## Tech

- Go
  - Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout?tab=readme-ov-file#go-directories)
- [go-chi/chi](https://github.com/go-chi/chi)
- [rs/zerolog](https://github.com/rs/zerolog)

## Discussions

### Thoughts on inlining templates & assets

This is such a good feature. But I did not use it for this project. And I'd expect even though it is incredibly awesome, I would not use it in most projects either.

While the single binary deployment experience is awesome, the overhead for copying two additional folders (`/templates` and `/public`) is negligible. In practice, it's adding two more lines to my Dockerfile. Having those assets as files also let's me edit the templates and let's me tweak the css without the binary recompilation. Sure, in this case, as a toy this is no big deal. Inlining causes a flavor of locking I'm not too fond of.

Meanwhile, in other cases I would recommend it. Are you using go proxy as an interstitial an idp workflow page? Great, keep that inlined because changing it should get explicitly versioned because it's too important to not get vetted by the formal process.
