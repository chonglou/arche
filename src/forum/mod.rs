pub mod users;
pub mod models;
pub mod dao;

use rocket;

pub fn mount(app: rocket::Rocket) -> rocket::Rocket {
    return app.mount("/users", routes![users::get_sign_in, users::get_sign_up]);
}
