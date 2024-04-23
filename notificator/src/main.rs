use notify_rust::Notification;
use tungstenite::{connect, Message};
use url::Url;

fn main() {
    env_logger::init();

    let (mut socket, response) =
        connect(Url::parse("ws://localhost:8000/ws?api-key=123").unwrap()).expect("Can't connect");

    println!("Connected to the server");
    println!("Response HTTP code: {}", response.status());
    println!("Response contains the following headers:");
    for (ref header, _value) in response.headers() {
        println!("* {}", header);
    }

    socket
        .send(Message::Text("Hello WebSocket".into()))
        .unwrap();
    loop {
        let msg = socket.read().expect("Error reading message");
        println!("Received: {}", msg);
        let _ = Notification::new()
            .summary("Firefox News")
            .body("This will almost look like a real firefox notification.")
            .icon("firefox")
            .show();
    }
    // socket.close(None);
}
