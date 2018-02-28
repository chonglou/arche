dist=dist
pkg=github.com/chonglou/arche/web
theme=moon

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`

build: www my
	cd $(dist) && mkdir -p tmp public
	cd $(dist) && tar cfJ ../$(dist).tar.xz *

www:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/arche main.go
	-cp -r db locales templates themes LICENSE README.md $(dist)/

my:
	cd dashboard && npm run build
	-cp -r dashboard/build $(dist)/dashboard

clean:
	-rm -r $(dist) $(dist).tar.xz dashboard/build

init:
	govendor sync
	cd dashboard && npm install
