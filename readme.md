# PuGo

a simple static site generator

### Install

    go get github.com/fuxiaohei/pugov1

go to source directory, run:

    go build main.go

then run the binary file:

    ./pugov1 version

print current version:

    0.1.0

### Generate

init default site:

    ./pugov1 init

build site contents to html pages:

    ./pugov1 build

it generates all html files in **_dest** directory.

preview the site:

    ./pugov1 server

then visit **http://localhost:9899/** to display the webpage in browser

full commands usages in [wiki-Commands](#).

### Deploy

the files in **_dest** are a whole site. Put them to http server (such as nginx) directory.