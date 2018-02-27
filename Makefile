dist=dist
pkg=github.com/chonglou/arche/web
theme=moon

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`

build: www api
	cd $(dist) && mkdir -p tmp public
	cd $(dist) && tar cfJ ../$(dist).tar.xz *

api:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/arche main.go
	-cp -r db locales templates LICENSE README.md $(dist)/

www:
	cd frontend && npm run build
	mkdir -p $(dist)/frontend
	-cp -r frontend/.next frontend/package.json frontend/package-lock.json $(dist)/frontend/

clean:
	-rm -r $(dist) $(dist).tar.xz frontend/.next

init:
	govendor sync
	cd frontend && npm install
