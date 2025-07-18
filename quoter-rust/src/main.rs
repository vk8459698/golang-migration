use serde::Serialize;
use warp::Filter;

#[derive(Serialize)]
struct PingResponse {
    status: &'static str,
    message: &'static str,
}

#[tokio::main]
async fn main() {
    let port = 8080;

    let ping = warp::path!("ping").and(warp::get()).map(|| {
        let response = PingResponse {
            status: "ok",
            message: "pong",
        };
        warp::reply::json(&response)
    });

    println!("quoter rust started on port {}", port);

    warp::serve(ping).run(([0, 0, 0, 0], port)).await;
}
