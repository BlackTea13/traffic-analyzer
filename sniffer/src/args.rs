use clap::Parser;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    // The port to listen for packets
    #[arg(short = 'p', long = "port", default_value = "9995")]
    pub port: u32,

    // The address of a Kafka Broker
    #[arg(short = 'k', long = "kafka", default_value = "localhost:19092")]
    pub kafka_broker: String,

    // The number of threads sniffing and processing
    #[arg(short = 'n', long = "num_threads", default_value = "3")]
    pub num_threads: u8,
}
