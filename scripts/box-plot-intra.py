import argparse
import threading
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import seaborn as sns
import json
import os
import glob

import logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


sent_times = {}
rtts_by_index = {}
timestamps_by_index = {}
servers_by_index = {}
alerts = []
lost_ids = set()
lock = threading.Lock()
stop_event = threading.Event()

def load_results(filepath):
    """
    Function to load results from a JSON file.
    :param filepath: Path to the JSON file containing results.
    :return: Dictionary with timestamps and RTTs.
    """

    with open(filepath, 'r') as f:
        data = json.load(f)

    sent_total = data.get('sent_total', 0)
    rtts_by_index = {int(k): v for k, v in data.get('rtts', {}).items()}
    timestamps_by_index = {int(k): v for k, v in data.get('timestamps', {}).items()}
    servers_by_index = {int(k): v for k, v in data.get('servers', {}).items()}
    alerts = data.get('alerts', [])
    
    return rtts_by_index, timestamps_by_index, servers_by_index, sent_total, alerts


def process_alert_phases(alerts):
    """
    Process alerts from a single run to extract phase durations for single domain migration.
    :param alerts: List of alerts from one run.
    :return: Dictionary mapping phase names to durations.
    """
    # Sort alerts by timestamp
    sorted_alerts = sorted(alerts, key=lambda x: x.get("timestamp", 0))
    
    # Create a lookup for alert timestamps by name
    # For duplicate alert names, use the last occurrence (highest timestamp)
    alert_times = {}
    for alert in sorted_alerts:
        name = alert.get("name", "")
        timestamp = alert.get("timestamp")
        if timestamp is not None:
            # Always update to use the latest timestamp for this alert name
            alert_times[name] = timestamp
    
    # Define phase boundaries for single domain migration between 2 clusters
    phase_definitions = {
        "New App\nInstantiation": ("node-inst-init", "node-inst-done"),
        "Network\nInterface\nSwitch": ("new-server-started", "new-server-switch"),
        "Old App\nTermination": ("old-server-signal-15", "old-server-killed"),
        "Migration\nTime": ("migration-init", "old-server-killed"),
    }
    
    phase_durations = {}
    
    # Calculate duration for each phase
    for phase_name, (start_alert, end_alert) in phase_definitions.items():
        if start_alert in alert_times and end_alert in alert_times:
            start_time = alert_times[start_alert]
            end_time = alert_times[end_alert]
            duration = end_time - start_time
            phase_durations[phase_name] = duration
        else:
            logging.warning(f"Missing alerts for phase {phase_name}: {start_alert} or {end_alert}")
    
    return phase_durations


def load_multiple_results(folder_path):
    """
    Function to load and aggregate results from multiple JSON files in a folder.
    :param folder_path: Path to the folder containing JSON result files.
    :return: Aggregated data from all runs.
    """
    if not os.path.exists(folder_path):
        logging.error(f"Folder {folder_path} does not exist")
        return {}, {}, {}, 0, {}
    
    # Find all JSON files in the folder
    json_files = glob.glob(os.path.join(folder_path, "*.json"))
    
    if not json_files:
        logging.warning(f"No JSON files found in folder {folder_path}")
        return {}, {}, {}, 0, {}
    
    logging.info(f"Found {len(json_files)} JSON files to process")
    
    # Aggregate data from all files
    all_rtts = {}
    all_timestamps = {}
    all_servers = {}
    total_sent = 0
    all_transitions = {}  # Dictionary mapping transition names to lists of durations
    
    current_index_offset = 0
    
    for file_path in sorted(json_files):
        try:
            logging.info(f"Processing file: {os.path.basename(file_path)}")
            rtts, timestamps, servers, sent_total, alerts = load_results(file_path)
            
            # Offset indices to avoid conflicts between runs
            for idx, rtt in rtts.items():
                all_rtts[idx + current_index_offset] = rtt
            
            for idx, timestamp in timestamps.items():
                all_timestamps[idx + current_index_offset] = timestamp
                
            for idx, server in servers.items():
                all_servers[idx + current_index_offset] = server
            
            # Process alerts for this run to get phase durations
            run_phases = process_alert_phases(alerts)
            
            # Aggregate phase durations across runs
            for phase_name, duration in run_phases.items():
                if phase_name not in all_transitions:
                    all_transitions[phase_name] = []
                all_transitions[phase_name].append(duration)
            
            total_sent += sent_total
            
            # Update offset for next file
            if rtts:
                current_index_offset += max(rtts.keys()) + 1
                
        except Exception as e:
            logging.error(f"Error processing file {file_path}: {e}")
            continue
    
    logging.info(f"Aggregated data from {len(json_files)} runs: {len(all_transitions)} unique phases")
    return all_rtts, all_timestamps, all_servers, total_sent, all_transitions



def plot_results(sent_total, phase_durations, output=None):
    """
    Function to create a box plot showing migration phase durations for single domain migration.
    :param rtts_by_index: Dictionary mapping message index to RTT.
    :param timestamps_by_index: Dictionary mapping message index to timestamp.
    :param servers_by_index: Dictionary mapping message index to server ID.
    :param sent_total: Total number of messages sent.
    :param phase_durations: Dictionary mapping phase names to lists of durations across runs.
    :param output: Optional output file to save the plot.
    """
    if not phase_durations:
        logging.warning("No phase durations found for analysis")
        return

    # Create box plot
    plt.figure(figsize=(18, 10))
    
    # Set larger font sizes
    plt.rcParams.update({'font.size': 14})
    
    # Prepare data for box plot
    box_data = []
    labels = []
    
    for phase_name, durations in phase_durations.items():
        if durations:
            box_data.append(durations)
            labels.append(phase_name)
    
    if not box_data:
        logging.warning("No valid phase duration data found for box plot")
        return
    
    # Create the box plot
    box_plot = plt.boxplot(box_data, labels=labels, patch_artist=True, showfliers=False)
    
    # Color the boxes with distinct colors for each phase
    colors = sns.color_palette('Set2', len(box_data))
    for patch, color in zip(box_plot['boxes'], colors):
        patch.set_facecolor(color)
        patch.set_alpha(0.7)
    
    plt.xlabel("Migration Phases", fontsize=28)
    plt.ylabel("Duration (seconds)", fontsize=28)
    plt.title("Single Domain Migration Phase Duration Distribution", fontsize=32)
    plt.xticks(rotation=0, ha='center', fontsize=26)
    plt.yticks(fontsize=26)
    plt.ylim(0, 80)
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    # Print phase statistics
    logging.info(f"\nSingle Domain Migration Phase Duration Statistics (from {sent_total} total messages):")
    for phase_name, durations in phase_durations.items():
        if durations:
            logging.info(f"{phase_name}:")
            logging.info(f"  Count: {len(durations)}")
            logging.info(f"  Min: {min(durations):.3f}s")
            logging.info(f"  Max: {max(durations):.3f}s")
            logging.info(f"  Mean: {sum(durations)/len(durations):.3f}s")
            if len(durations) > 1:
                import statistics
                logging.info(f"  Median: {statistics.median(durations):.3f}s")
                logging.info(f"  StdDev: {statistics.stdev(durations):.3f}s")

    if output:
        plt.savefig(output)
    else:
        try:
            plt.show()
        except KeyboardInterrupt:
            logging.info("Plot window closed via Ctrl+C.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Plot single domain migration phase durations from multiple runs in a folder")
    parser.add_argument("-d", "--folder", help="Path to the folder containing JSON result files", default="node-teste/")
    parser.add_argument("-o", "--output", help="Output file to save plot (default: box-plot-intra.pdf)", default="box-plot-intra.pdf")

    args = parser.parse_args()
    rtts_by_index, timestamps_by_index, servers_by_index, sent_total, phase_durations = load_multiple_results(args.folder)

    # plot the results after processing all runs
    plot_results(sent_total, phase_durations, args.output)