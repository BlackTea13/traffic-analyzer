use crate::kafka::produce_packets;
use crate::packet::Packet;
use bytes::Bytes;
use netflow_parser::{NetflowPacketResult, NetflowParser};
use rdkafka::producer::FutureProducer;
use std::sync::Arc;
use tokio::net::UdpSocket;
use tokio::sync::mpsc;

pub struct Sniffer {
    socket: Arc<UdpSocket>,
    num_threads: u8,
    processors: Vec<mpsc::Sender<Bytes>>,
}

impl Sniffer {
    pub async fn new(
        port: u32,
        kafka_producer: FutureProducer,
        num_threads: u8,
    ) -> Result<Self, anyhow::Error> {
        let address = format!("0.0.0.0:{}", port);
        let socket = UdpSocket::bind(address).await?;

        let producer = Arc::new(kafka_producer);

        let mut sniffer = Sniffer {
            socket: Arc::new(socket),
            num_threads,
            processors: Vec::new(),
        };

        for _ in 0..num_threads {
            let (sender, receiver) = mpsc::channel(1024);
            sniffer.processors.push(sender);

            tokio::task::spawn(Sniffer::process(receiver, producer.clone()));
        }

        Ok(sniffer)
    }

    /// Spawns sniffers
    pub async fn start(&self) -> Result<(), anyhow::Error> {
        for i in 0..self.num_threads {
            let sender = self.processors[i as usize].clone();
            tokio::task::spawn(Sniffer::sniff(self.socket.clone(), sender));
        }
        Ok(())
    }

    // Sniffs packets and sends them to the processor
    async fn sniff(socket: Arc<UdpSocket>, sender: mpsc::Sender<Bytes>) {
        let mut buf = vec![0; 65536];

        loop {
            match socket.recv_from(&mut buf).await {
                Ok((len, _addr)) => {
                    let packet = buf[..len].to_vec();
                    let _ = sender.send(Bytes::from(packet)).await;
                }
                Err(e) => {
                    eprintln!("Error reading from socket: {}", e)
                }
            }
        }
    }

    /// Parses and processes packets
    async fn process(mut receiver: mpsc::Receiver<Bytes>, producer: Arc<FutureProducer>) {
        let mut parser = NetflowParser::default();
        let producer = producer.clone();

        while let Some(packet) = receiver.recv().await {
            let parsed = parser.parse_bytes(&packet);

            for flow_result in parsed {
                match flow_result {
                    NetflowPacketResult::V5(v5) => {
                        let packet = Packet::from(v5);
                        let _ = produce_packets(producer.clone(), packet).await;
                    }
                    NetflowPacketResult::IPFix(ipfix) => {
                        let packet = if let Ok(packet) = Packet::try_from(ipfix) {
                            packet
                        } else {
                            eprintln!("Could not parse IPFIX packet");
                            continue;
                        };
                        let _ = produce_packets(producer.clone(), packet).await;
                    }
                    _ => {
                        eprintln!("Unsupported packet version")
                    }
                }
            }
        }
    }
}
