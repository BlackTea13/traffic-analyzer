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

## Flux Queries for Grafana Dashboards

```
// Packet Count by Source Country
from(bucket: "bucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "packet")
  |> filter(fn: (r) => r["_field"] == "size")
  |> group(columns: ["source_country"])
  |> filter(fn: (r) => r["source_country"] != "Value")
  |> aggregateWindow(every: 5h, fn: count, createEmpty: false)
  |> yield()

// Mean Packet Size by Source Country
from(bucket: "bucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "packet")
  |> filter(fn: (r) => r["_field"] == "size")
  |> group(columns: ["source_country"])
  |> filter(fn: (r) => r["source_country"] != "null")
  |> aggregateWindow(every: 1d, fn: mean)
  |> sort(columns: ["_time"])
  |> filter(fn: (r) => r._value >= 5000)
```
