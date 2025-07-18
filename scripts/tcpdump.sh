#!/bin/bash

# Start capture
tcpdump -i any -s 0 -w api_flow.pcap host your-api-server.com and port 443 &
TCPDUMP_PID=$!

# Run your API test here
echo "Run your API requests now..."
sleep 30  # Adjust based on your test duration

# Stop capture
kill $TCPDUMP_PID

# Analysis
echo "=== PACKET ANALYSIS ==="
TOTAL_PACKETS=$(tcpdump -r api_flow.pcap 2>/dev/null | wc -l)
echo "Total packets: $TOTAL_PACKETS"

# Calculate total bytes
TOTAL_BYTES=$(tcpdump -r api_flow.pcap -q -n 2>/dev/null | grep -oE 'length [0-9]+' | awk '{sum += $2} END {print sum}')
echo "Total bytes: $TOTAL_BYTES"

# Connection overhead
TCP_HANDSHAKE=$(tcpdump -r api_flow.pcap -n 'tcp[tcpflags] & tcp-syn != 0' 2>/dev/null | wc -l)
echo "TCP handshake packets: $TCP_HANDSHAKE"

# Average packet size
if [ $TOTAL_PACKETS -gt 0 ]; then
    AVG_SIZE=$((TOTAL_BYTES / TOTAL_PACKETS))
    echo "Average packet size: $AVG_SIZE bytes"
fi