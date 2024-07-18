use crate::packet::Packet;
use crate::TOPIC;
use rdkafka::producer::{FutureProducer, FutureRecord};
use serde_json::Result;
use std::sync::Arc;
use tokio::time::Duration;

pub async fn produce_packets(producer: Arc<FutureProducer>, packet: Packet) -> Result<()> {
    let producer: &FutureProducer = &producer.clone();
    let futures = packet
        .flows
        .into_iter()
        .map(|f| f.to_json().to_string())
        .map(|data| async move {
            let data: String = data;
            let record: FutureRecord<String, String> = FutureRecord::to(TOPIC).payload(&data);
            let _ = producer.send(record, Duration::from_secs(0)).await;
        })
        .collect::<Vec<_>>();

    for future in futures {
        future.await;
    }

    Ok(())
}
