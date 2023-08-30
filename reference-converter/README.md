# Reference Converter

This program converts the NGINX reference documentation from it's XML schema to JSON. The generated JSON is available as an npm package in the reference-lib folder and can be used for static content generation, markdoc tags, monaco plugins, etc.

## Design

```mermaid
flowchart
    fetch_atom[read latest version from atom feed]
    download_tarball[download tarball of all XML]
    parse_xml[parse XML]
    render_md[translate prose to markdown]
    write[write JSON to disk]
    done((done))

    fetch_atom --> download_tarball --> parse_xml --> render_md --> write --> done
```

The NGINX docs are publicly available at <http://hg.nginx.org/nginx.org>, in XML that's a mix of data and prose (`<para>` tags contain markup). The `<para>` contents will be translated in-order to generate equivalent markdown.

The atom feed at <http://hg.nginx.org/nginx.org/atom-log> will tell us if there is updated content.

A scheduled github pipeline ensures that we have up-to-date reference information.

```mermaid
flowchart
    run[./reference-converter]
    diff{json file has changed?}
    open[open a PR with the changes]
    slack[send slack notification]
    done((done))

    run --> diff -->|N| done
    diff -->|Y| open --> slack --> done
    run -->|errored out| slack
```

## Usage

```bash
make devtools-image
make build
./dist/reference-converter --dst <output-path>
```
