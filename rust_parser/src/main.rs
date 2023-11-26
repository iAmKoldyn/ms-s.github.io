use csv::ReaderBuilder;
use lapin::{
    options::BasicPublishOptions,
    BasicProperties, Connection, ConnectionProperties,
};
use log::{info, warn};
use reqwest;
use scraper::{Html, Selector};
use serde::{Deserialize, Serialize};
use std::{error::Error, fs::File, env};

#[derive(Debug, Deserialize, Serialize)]
struct User {
    url: String,
    name: String,
    html: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    env::set_var("RUST_LOG", "info");
    env_logger::init();

    let amqp_addr = "amqp://guest:guest@localhost:5672/%2f";
    let conn = Connection::connect(&amqp_addr, ConnectionProperties::default()).await?;
    println!("Connected to RabbitMQ at {}", amqp_addr);
    let channel = conn.create_channel().await?;
    let queue = "users_from_habr";

    let mut users = vec![];
    let file = File::open("./habr.csv")?;
    let mut rdr = ReaderBuilder::new().delimiter(b',').from_reader(file);
    for result in rdr.records() {
        let record = result?;
        let url = &record[0];
        info!("URL: {}", url);
        users.push(User {
            url: url.to_string(),
            name: String::new(),
            html: String::new(),
        });
    }

    for user in &mut users {
        match reqwest::get(&user.url).await {
            Ok(response) => {
                let html_content = response.text().await?;
                let document = Html::parse_document(&html_content);
                let selector = Selector::parse(".page-title__title").unwrap();
                user.name = document.select(&selector).next().map_or_else(|| String::new(), |n| n.inner_html());
                user.html = html_content;
                info!("Name: {}", user.name);

                let message = serde_json::to_string(&user)?;
                channel.basic_publish("", queue, BasicPublishOptions::default(), message.into_bytes(), BasicProperties::default()).await?;
                println!("Published message for URL: {}", user.url);
            },
            Err(e) => warn!("Error fetching URL {}: {}", user.url, e),
        };
    }

    Ok(())
}
