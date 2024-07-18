use anyhow::anyhow;
use chrono::Utc;
use netflow_parser::static_versions::v5::V5;
use netflow_parser::variable_versions::common::DataNumber;
use netflow_parser::variable_versions::common::FieldValue;
use netflow_parser::variable_versions::ipfix::IPFix;
use netflow_parser::variable_versions::ipfix_lookup::IPFixField;
use serde::Serialize;
use serde_json::{json, Value};

#[derive(Debug)]
pub struct Packet {
    pub flows: Vec<Flow>,
}

#[derive(Debug, Serialize)]
pub struct Flow {
    src_addr: String,
    src_port: String,
    dst_addr: String,
    dst_port: String,
    packet_size: u32,
    timestamp: String,
}

impl Flow {
    pub fn to_json(&self) -> Value {
        json!({
            "SrcIp": self.src_addr,
            "SrcPort": self.src_port,
            "DestIp": self.dst_addr,
            "DestPort": self.dst_port,
            "Size": self.packet_size,
            "Timestamp": self.timestamp
        })
    }
}

impl From<V5> for Packet {
    fn from(v5_packet: V5) -> Self {
        let mut flows = Vec::new();
        for flowset in v5_packet.flowsets {
            let packet = Flow {
                src_addr: flowset.src_addr.to_string(),
                src_port: flowset.src_port.to_string(),
                dst_addr: flowset.dst_addr.to_string(),
                dst_port: flowset.dst_port.to_string(),
                packet_size: flowset.d_octets as u32,
                timestamp: Utc::now().to_string(),
            };

            flows.push(packet);
        }

        Packet { flows }
    }
}

impl TryFrom<IPFix> for Packet {
    type Error = anyhow::Error;

    fn try_from(ipfix_packet: IPFix) -> Result<Self, Self::Error> {
        let mut flows: Vec<Flow> = Vec::new();
        for flowset in ipfix_packet.flowsets {
            let data = flowset
                .body
                .data
                .ok_or_else(|| anyhow!(format!("Could not retrieve flowset data")))?;

            for data_field in data.data_fields {
                let mut src_addr = String::new();
                let mut src_port = String::new();
                let mut dst_addr = String::new();
                let mut dst_port = String::new();
                let mut packet_size = 0;
                let timestamp = Utc::now();

                for (_, (field, value)) in data_field {
                    match field {
                        IPFixField::SourceIpv4address => {
                            if let FieldValue::Ip4Addr(addr) = value {
                                src_addr = addr.to_string();
                            } else {
                                return Err(anyhow!("Only IPv4 addresses are supported"));
                            }
                        }
                        IPFixField::DestinationIpv4address => {
                            if let FieldValue::Ip4Addr(addr) = value {
                                dst_addr = addr.to_string();
                            } else {
                                return Err(anyhow!("Only IPv4 addresses are supported"));
                            }
                        }
                        IPFixField::SourceTransportPort => {
                            if let FieldValue::DataNumber(data_number) = value {
                                src_port = match data_number {
                                    DataNumber::U8(port) => port.to_string(),
                                    DataNumber::U16(port) => port.to_string(),
                                    DataNumber::U24(port) => (port as u16).to_string(),
                                    DataNumber::U32(port) => (port as u32).to_string(),
                                    _ => return Err(anyhow!("Data number not supported")),
                                };
                            }
                        }
                        IPFixField::DestinationTransportPort => {
                            if let FieldValue::DataNumber(data_number) = value {
                                dst_port = match data_number {
                                    DataNumber::U8(port) => port.to_string(),
                                    DataNumber::U16(port) => port.to_string(),
                                    DataNumber::U24(port) => (port as u16).to_string(),
                                    DataNumber::U32(port) => (port as u32).to_string(),
                                    _ => {
                                        return Err(anyhow!("Data number not supported"));
                                    }
                                };
                            }
                        }
                        IPFixField::PacketTotalCount => {
                            if let FieldValue::DataNumber(data_number) = value {
                                packet_size = match data_number {
                                    DataNumber::U8(value) => value as i32,
                                    DataNumber::U16(value) => value as i32,
                                    DataNumber::U24(value) => value as i32,
                                    DataNumber::U32(value) => value as i32,
                                    DataNumber::I32(value) => value,
                                    _ => {
                                        return Err(anyhow!("Data number not supported"));
                                    }
                                }
                            }
                        }
                        _ => {}
                    }
                }

                let flow = Flow {
                    src_addr,
                    src_port,
                    dst_addr,
                    dst_port,
                    packet_size: packet_size as u32,
                    timestamp: timestamp.to_string(),
                };

                flows.push(flow);
            }
        }

        Ok(Packet { flows })
    }
}
