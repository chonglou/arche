# ARCHE

Blog Tool, Publishing Platform, and CMS (By Go and React).

## Usage

- Install go

  ```bash
  zsh < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  gvm install go1.10 -B
  gvm use go1.10 --default
  ```

- npm

  ```bash
  mkdir ~/.npm-global
  npm config set prefix '~/.npm-global'
  export PATH=~/.npm-global/bin:$PATH
  ```

- Build

  ```bash
  go get -u github.com/chonglou/arche
  cd $GOPATH/src/github.com/chonglou/arche
  make init
  make clean
  make
  ```

- create database(for postgresql)

  ```sql
  CREATE DATABASE db-name WITH ENCODING = 'UTF8';
  CREATE USER user-name WITH PASSWORD 'change-me';
  GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
  ```

- Generate a random key

  ```bash
  openssl rand -base64 32
  ```

## Atom plugins

enable autosave

- go-plus
- file-icons
- atom-beautify(enable newline, beautify on save; need python-sqlparse)
- language-babel
- language-ini

## Documents

- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)

- [favicon.ico](http://icoconvert.com/)

- [smver](http://semver.org/)

- [banner.txt](http://patorjk.com/software/taag/)

- [AWS](http://docs.aws.amazon.com/general/latest/gr/rande.html)

- [Bootstrap](http://getbootstrap.com/)

- [Ant Design](https://ant.design/docs/react/introduce)

- [Ant Design Pro](https://pro.ant.design/docs/getting-started)

- [Font Awesome](https://fontawesome.com/how-to-use/js-component-packages)
