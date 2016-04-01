# fessboxconfig
Configuration manager for fessbox

# checklist

- [WIP]  JSON API to work with the files
- [x] Asterisk configuration format
  - [x] context
  - [x] options
  - [ ] templates
  - [ ] include
  - [ ] exec
  - [x] comments
    - [x] single line comments
    - [x] block comments
- [WIP] Device autodetection



# Links
- [asterisk](http://www.asterisk.org/)
- [asterisk configuration format](https://wiki.asterisk.org/wiki/display/AST/Config+File+Format)


# Developing

You need a working Golang environment

First install
```bash
go get github.com/FarmRadioHangar/fessboxconfig/...
```


Charge to the root of the installed repository

```bash
cd $GOPATH/github.com/FarmRadioHangar/fessboxconfig
```

If you have asterisk installed and want to run in production mode

```bash
make start
```

If you just want to test in dev mode( No asterisk installation is required)

```bash
make dev
```
