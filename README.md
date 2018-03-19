# arche

Blog Tool, Publishing Platform, and CMS (By Rust and React).

## Install

-   install rust


    curl https://sh.rustup.rs -sSf | sh
    rustup default nightly
    cargo install rustfmt-nightly
    cargo install racer
    rustup component add rust-src

-   upgrade

    rustup update
    cargo update

-   add to your .zshrc

    export PATH="$HOME/.cargo/bin:$PATH"
    export RUST_SRC_PATH="$(rustc --print sysroot)/lib/rustlib/src/rust/src"

-   test racer

    racer complete std::io::B

-   test run

    cargo run -- --version

## Atom plugins

enable autosave

-   language-rust
-   racer
-   file-icons
-   atom-beautify(enable newline, beautify on save; need python-sqlparse)
-   language-babel
-   language-ini

## Notes

-   Generate a random key

    openssl rand -base64 32

-   ~/.npmrc

    prefix=${home}/.npm-packages

-   Create database

    CREATE DATABASE db-name WITH ENCODING = 'UTF8';
    CREATE USER user-name WITH PASSWORD 'change-me';
    GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;

## Documents

-   [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)

-   [favicon.ico](http://icoconvert.com/)

-   [smver](http://semver.org/)

-   [banner.txt](http://patorjk.com/software/taag/)

-   [The Rust Programming Language](https://doc.rust-lang.org/book/second-edition/)
