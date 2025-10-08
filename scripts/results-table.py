#!/usr/bin/env python3
import json
import glob
from datetime import datetime
from collections import defaultdict
import statistics

def parse_timestamp(ts):
    """Parse ISO format timestamp to datetime object."""
    return datetime.fromisoformat(ts)

def extract_flow_type(event_name):
    """Extract flow type from event name (e.g., 'oo-init-create-federation' -> 'create-federation')."""
    parts = event_name.split('-')
    # Remove prefixes like 'oo', 'po', 'federation-po' and suffixes like 'init', 'done'
    if 'init' in parts:
        idx = parts.index('init')
        return '-'.join(parts[idx+1:])
    elif 'done' in parts:
        idx = parts.index('done')
        return '-'.join(parts[idx+1:])
    return event_name

def analyze_flows(json_files):
    """Analyze flow durations from JSON event files."""
    # Group events by flow_id and file
    flows_by_type = defaultdict(list)

    for json_file in json_files:
        with open(json_file, 'r') as f:
            events = json.load(f)

        # Group events by flow_id
        flows = defaultdict(list)
        for event in events:
            flows[event['flow_id']].append(event)

        # Calculate duration for each flow
        for flow_id, flow_events in flows.items():
            if len(flow_events) < 2:
                continue

            # Sort events by timestamp
            flow_events.sort(key=lambda e: parse_timestamp(e['timestamp']))

            first_event = flow_events[0]
            last_event = flow_events[-1]

            # Extract flow type from event name
            flow_type = extract_flow_type(first_event['name'])

            # Calculate duration in milliseconds
            start = parse_timestamp(first_event['timestamp'])
            end = parse_timestamp(last_event['timestamp'])
            duration_ms = (end - start).total_seconds() * 1000

            flows_by_type[flow_type].append(duration_ms)

    return flows_by_type

def print_statistics(flows_by_type):
    """Print statistics for each flow type."""
    print(f"{'Flow Type':<40} {'Count':<8} {'Min (ms)':<12} {'Max (ms)':<12} {'Avg (ms)':<12} {'Median (ms)':<12} {'Std Dev (ms)':<12}")
    print("=" * 130)

    for flow_type in sorted(flows_by_type.keys()):
        durations = flows_by_type[flow_type]

        count = len(durations)
        min_val = min(durations)
        max_val = max(durations)
        avg_val = statistics.mean(durations)
        median_val = statistics.median(durations)
        std_dev = statistics.stdev(durations) if count > 1 else 0.0

        print(f"{flow_type:<40} {count:<8} {min_val:<12.2f} {max_val:<12.2f} {avg_val:<12.2f} {median_val:<12.2f} {std_dev:<12.2f}")

def main():
    # Find all JSON files in the results folder
    json_files = glob.glob('results/*.json')

    if not json_files:
        print("No JSON files found in scripts/results/")
        return

    print(f"Analyzing {len(json_files)} JSON file(s)...\n")

    flows_by_type = analyze_flows(json_files)
    print_statistics(flows_by_type)

if __name__ == '__main__':
    main()
