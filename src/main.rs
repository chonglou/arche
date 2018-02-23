extern crate arche;
extern crate env_logger;

fn main() {
    env_logger::init();
    match arche::console::run() {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
}
