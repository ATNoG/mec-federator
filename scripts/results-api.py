from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Any, List
from datetime import datetime
import uvicorn
import json
from pathlib import Path

app = FastAPI()

# Event message model matching the Go ResultsMessage struct
class EventMessage(BaseModel):
    name: str
    message: str
    value: Any

# Stored event with timestamp and flow tracking
class StoredEvent(BaseModel):
    name: str
    message: str
    value: Any
    timestamp: datetime
    flow_id: int

# File storage configuration
EVENTS_FILE = Path("events.json")

# In-memory storage
events: List[StoredEvent] = []
current_flow_id: int = 0

def save_events():
    """Save events to file"""
    with open(EVENTS_FILE, "w") as f:
        json.dump([e.model_dump(mode="json") for e in events], f, indent=2, default=str)

def load_events():
    """Load events from file"""
    global current_flow_id
    if EVENTS_FILE.exists():
        with open(EVENTS_FILE, "r") as f:
            data = json.load(f)
            loaded_events = [StoredEvent(**e) for e in data]
            if loaded_events:
                current_flow_id = max(e.flow_id for e in loaded_events)
            return loaded_events
    return []

# Load existing events on startup
events = load_events()

@app.post("/alert")
async def receive_event(event: EventMessage):
    """Receive and store an event message"""
    global current_flow_id

    # If message is "reset", start a new flow
    print(event.message)
    if event.message == "reset":
        current_flow_id += 1

    stored_event = StoredEvent(
        name=event.name,
        message=event.message,
        value=event.value,
        timestamp=datetime.now(),
        flow_id=current_flow_id
    )
    events.append(stored_event)
    save_events()
    return {"status": "success", "event_id": len(events) - 1, "flow_id": current_flow_id}

@app.get("/events")
async def get_events(name: str = None, flow_id: int = None):
    """Retrieve all stored events, optionally filtered by name and/or flow_id"""
    filtered = events
    if name:
        filtered = [e for e in filtered if e.name == name]
    if flow_id is not None:
        filtered = [e for e in filtered if e.flow_id == flow_id]
    return {"count": len(filtered), "events": filtered}

@app.get("/flows")
async def get_flows():
    """Get all flows with their events"""
    flows = {}
    for event in events:
        if event.flow_id not in flows:
            flows[event.flow_id] = []
        flows[event.flow_id].append(event)

    return {
        "total_flows": len(flows),
        "current_flow_id": current_flow_id,
        "flows": {fid: {"count": len(evts), "events": evts} for fid, evts in flows.items()}
    }

@app.delete("/events")
async def clear_events():
    """Clear all stored events"""
    global current_flow_id
    events.clear()
    current_flow_id = 0
    save_events()
    return {"status": "success", "message": "All events cleared"}

@app.get("/flows/{flow_id}/timing")
async def get_flow_timing(flow_id: int):
    """Get timing analysis for a specific flow"""
    flow_events = [e for e in events if e.flow_id == flow_id]

    if not flow_events:
        raise HTTPException(status_code=404, detail="Flow not found")

    if len(flow_events) < 2:
        return {"flow_id": flow_id, "message": "Not enough events for timing analysis"}

    # Sort by timestamp
    flow_events.sort(key=lambda e: e.timestamp)

    timings = []
    for i in range(1, len(flow_events)):
        duration = (flow_events[i].timestamp - flow_events[i-1].timestamp).total_seconds()
        timings.append({
            "from": flow_events[i-1].message,
            "to": flow_events[i].message,
            "duration_seconds": duration,
            "timestamp_start": flow_events[i-1].timestamp,
            "timestamp_end": flow_events[i].timestamp
        })

    total_duration = (flow_events[-1].timestamp - flow_events[0].timestamp).total_seconds()

    return {
        "flow_id": flow_id,
        "total_events": len(flow_events),
        "total_duration_seconds": total_duration,
        "timings": timings
    }

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
