import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

# Generate sample data (replace with your actual data)
np.random.seed(42)

# Categories for x-axis
categories = ['Metrics\nRedistribution', 'Migration\nDecision', 'Target Pod\nInitialization', 
              'Target Pod\nReady', 'Source Pod\nTermination', 'Migration\nCompletion in\nOSM']

# Generate sample data for two scenarios
# In reality, you would load your actual experimental data
baseline_data = []
poc_data = []

# Sample data generation (replace with your actual measurements)
for i, cat in enumerate(categories):
    if i < 2:  # First two categories have very small values
        baseline_data.append(np.random.normal(0.01, 0.005, 50))
        poc_data.append(np.random.normal(0.008, 0.003, 50))
    elif i < 4:  # Middle categories
        baseline_data.append(np.random.normal(3.5, 0.5, 50))
        poc_data.append(np.random.normal(3.0, 0.4, 50))
    else:  # Last two categories
        baseline_data.append(np.random.normal(8.0, 1.0, 50))
        poc_data.append(np.random.normal(7.5, 0.8, 50))

# Create the plot
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 8), 
                               gridspec_kw={'height_ratios': [4, 1]}, 
                               sharex=True)

# Main plot (upper)
positions = np.arange(len(categories))
width = 0.35

# Create box plots
bp1 = ax1.boxplot(baseline_data, positions=positions - width/2, widths=width,
                   patch_artist=True, boxprops=dict(facecolor='lightblue'),
                   medianprops=dict(color='black', linewidth=1.5),
                   whiskerprops=dict(color='black'),
                   capprops=dict(color='black'),
                   flierprops=dict(marker='o', markersize=4, markerfacecolor='gray'))

bp2 = ax1.boxplot(poc_data, positions=positions + width/2, widths=width,
                   patch_artist=True, boxprops=dict(facecolor='lightcoral'),
                   medianprops=dict(color='black', linewidth=1.5),
                   whiskerprops=dict(color='black'),
                   capprops=dict(color='black'),
                   flierprops=dict(marker='o', markersize=4, markerfacecolor='gray'))

# Set labels and formatting for main plot
ax1.set_ylabel('Time (s)', fontsize=12)
ax1.set_ylim(0, 30)
ax1.grid(True, axis='y', alpha=0.3)
ax1.legend([bp1["boxes"][0], bp2["boxes"][0]], ['Baseline Scenario', 'PoC Scenario'], 
           loc='upper left', fontsize=10)

# Zoomed plot (lower) for small values
bp3 = ax2.boxplot(baseline_data[:2], positions=positions[:2] - width/2, widths=width,
                   patch_artist=True, boxprops=dict(facecolor='lightblue'),
                   medianprops=dict(color='black', linewidth=1.5))

bp4 = ax2.boxplot(poc_data[:2], positions=positions[:2] + width/2, widths=width,
                   patch_artist=True, boxprops=dict(facecolor='lightcoral'),
                   medianprops=dict(color='black', linewidth=1.5))

# Set labels and formatting for zoomed plot
ax2.set_ylim(0, 0.04)
ax2.set_ylabel('Time (s)', fontsize=10)
ax2.grid(True, axis='y', alpha=0.3)

# Set x-axis labels
ax2.set_xticks(positions)
ax2.set_xticklabels(categories, fontsize=10)

# Add title
fig.suptitle('Migration Stages Performance Comparison', fontsize=14, y=0.98)

# Adjust layout
plt.tight_layout()
plt.show()

