# Metro Vancouver Library Search

A small CLI tool that searches the catalogs of Metro Vancouver public libraries
at once and prints a unified summary of the results.

For each library it shows up to N matches with title, author, format, copies
available / total, a short description, and a direct catalog link.

## Build

The project is written in Go, and requires Go 1.20 or later to build:

```sh
go build -o metrovanlibsearch .
```

The build produces a single static binary with no runtime dependencies.

## Usage

```sh
./metrovanlibsearch "house of leaves"

# limit results to 5 per library, default is 3
./metrovanlibsearch --limit 5 "atwood"

# emit JSON instead of human-readable text
./metrovanlibsearch --json "the bear"

# filter by format code (e.g. only ebooks)
./metrovanlibsearch --format EBOOK "atwood"
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
