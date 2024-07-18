# Traffic Analyzer

Sniff Netflow V5/IPFix packets with the `sniffer`.

Enrich the packets with `mansamusa`.

Send the packets to DB with `processor`.

View beautiful graphs in Grafana.


## Sniffer

Run the sniffer by going to the `sniffer/` directory and run with 
```
cargo run
```

Checkout the `args.rs` file for args.
- `-p --port` for port
- `-k --kafka_broker` for the kafka address
- `-n --num_threads` for how many threads to run
