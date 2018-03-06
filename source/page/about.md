```toml
title = "About PuGo"
slug = ""
desc = "some words about pugo"
date = "2018-03-01 00:00"
update_date = "2018-03-06 22:43"
author = "pugo"
hover = "about"
template = "page.html"
lang = ""
sort = 0
draft = false

[meta]
  metadata = "this is meta data"
```

### Introduction

`PuGo` is a simple static site generator by [Golang](https://golang.org). It compiles [markdown](https://help.github.com/articles/markdown-basics/) to site pages with beautiful theme. No dependencies, cross platform and very fast.



![golang](/golang.png)




### Quick start

1. Download from [PuGo](http://pugo.io) and extract zip archive.
2. Run `pugo new site` to create new default site.
2. Run `pugo server`, open `http://localhost:9899` to visit.


### Commands

Run a command when run `pugo` executable file:

- `pugo new` create new site, post or page.
- `pugo build` build static files.
- `pugo server` build and serve static files.

More details in [Documentation](http://pugo.io/en/docs.html).

### Writing

`PuGo` support two kinds of content, `post` and `page`. you can create any `.md` file in proper directories in `source` directory. Read the [wiki](http://pugo.io/en/guide/write-new-post.html) to learn the layout format and more details.

### Compile

After you change your source `.md` files, run

    pugo build

To build static files.

### Customize

- Theme: theme documentation is [Here](http://pugo.io/en/docs/tpl/syntax.html)