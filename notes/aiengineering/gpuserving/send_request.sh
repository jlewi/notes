#!/bin/bash

curl -X POST \
  -H 'Content-Type: application/json' \
  -d $'{    
    "input": {
      "nlq": "EMISSING slowest traces",
      "cols": "[\'sli.latency\', \'duration_ms\', \'net.transport\', \'http.method\', \'error\', \'http.target\', \'http.route\', \'rpc.method\', \'ip\', \'http.request_content_length\', \'rpc.service\', \'apdex\', \'name\', \'message.type\', \'http.host\', \'service.name\', \'rpc.system\', \'http.scheme\', \'sli.platform-time\', \'type\', \'http.flavor\', \'span.kind\', \'dc.platform-time\', \'library.version\', \'status_code\', \'net.host.port\', \'net.host.ip\', \'app.request_id\', \'bucket_duration_ms\', \'library.name\', \'sli_product\', \'message.uncompressed_size\', \'rpc.grpc.status_code\', \'net.peer.port\', \'log10_duration_ms\', \'http.status_code\', \'status_message\', \'http.user_agent\', \'net.host.name\', \'span.num_links\', \'message.id\', \'parent_name\', \'app.cart_total\', \'num_products\', \'product_availability\', \'revenue_at_risk\', \'trace.trace_id\', \'trace.span_id\', \'ingest_timestamp\', \'http.server_name\', \'trace.parent_id\']"
    }
  }' \
  http://localhost:5000/predictions
