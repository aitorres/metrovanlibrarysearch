# Metro Vancouver Library Search

A small web and CLI tool that searches the catalogs of Metro Vancouver public
libraries at once and prints a unified summary of the results.

For each library it shows up to N matches with title, author, format, copies
available / total, a short description, and a direct catalog link.

## Build

The project is written in Go, and requires Go 1.20 or later to build:

```sh
go build -o metrovanlibsearch .
```

The build produces a single static binary with no runtime dependencies.

## Usage

The binary has two subcommands: `query` for a one-off CLI search and `serve`
for the local web UI.

```sh
./metrovanlibsearch query "house of leaves"

# limit results to 5 per library, default is 3
./metrovanlibsearch query --limit 5 "atwood"

# emit JSON instead of human-readable text
./metrovanlibsearch query --json "the bear"

# filter by format code (e.g. only ebooks)
./metrovanlibsearch query --format EBOOK "atwood"
```

### Web UI

Run a local web server with a minimal search page:

```sh
./metrovanlibsearch serve

# override the listen address (default :8080)
./metrovanlibsearch serve --addr :9090
```

The page shows a single search bar, you can submit queries
to get results from all libraries in a unified view.

### Docker

A `Dockerfile` is provided that builds a small static image and runs the web UI by default:

```sh
docker build -t metrovanlibsearch .
docker run --rm -p 8080:8080 metrovanlibsearch
```

### Format codes

The `--format` flag accepts a BiblioCommons format code and is passed
through to the catalog as a server-side filter. Codes are case-insensitive.
Common values include:

| Code            | Meaning             |
| --------------- | ------------------- |
| `BK`            | Book                |
| `EBOOK`         | eBook               |
| `AB`            | Audiobook (CD)      |
| `EAUDIO`        | Audiobook (digital) |
| `DVD`           | DVD                 |
| `BLU_RAY`       | Blu-ray             |
| `MUSIC_CD`      | Music CD            |
| `MUSIC_ONLINE`  | Streaming music     |
| `VIDEO_ONLINE`  | Streaming video     |
| `COMIC_BK`      | Comic book          |
| `GRAPHIC_NOVEL` | Graphic novel       |
| `MAG`           | Magazine            |
| `EMAG`          | eMagazine           |

## Supported libraries

The following libraries are supported via their Bibliocommons catalogs:

- Vancouver Public Library
- Burnaby Public Library
- Surrey Libraries
- New Westminster Public Library
- Richmond Public Library
- North Vancouver City Library
- West Vancouver Memorial Library
- Coquitlam Public Library
- Port Moody Public Library

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
