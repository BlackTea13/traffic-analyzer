use args::Args;
use clap::Parser;
use rdkafka::{producer::FutureProducer, ClientConfig};
use sniff::Sniffer;
use tokio::signal;

mod args;
mod kafka;
mod packet;
mod sniff;

pub const TOPIC: &str = "sniffed-packets";

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();

    println!("Starting Sniffer");
    println!("Connecting to Kafka at: {}", args.kafka_broker);

    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", args.kafka_broker)
        .set("message.timeout.ms", "5000")
        .create()
        .expect("Producer creation error");

    println!(
        "Listening on Port: {} with {} Sniffers",
        args.port, args.num_threads
    );
    let sniffer = Sniffer::new(args.port, producer, args.num_threads).await?;
    sniffer.start().await?;

    signal::ctrl_c().await?;
    Ok(())
}
