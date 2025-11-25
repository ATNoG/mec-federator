from kafka import KafkaConsumer, KafkaProducer
import sys
import os
import json
import time
import uuid
import yaml
import time

# Kafka configuration
bootstrap_servers = ['10.255.41.64:31999']
username = os.getenv('KAFKA_USERNAME', 'user1')
password = os.getenv('KAFKA_PASSWORD', 'password')
sasl_mechanism = 'PLAIN'
security_protocol = 'SASL_PLAINTEXT'

# Topic names
INSTANTIATE_TOPIC = 'federation_new_appi'
ONBOARD_TOPIC = 'federation_new_app'
TERMINATE_TOPIC = 'federation_remove_appi'
REMOVE_TOPIC = 'federation_remove_app'
RESPONSE_TOPIC = 'responses'

# Test data
federation_context_id = "9bc6a233-f990-4105-88bb-528d2479a37e"
app_pkg_id = "68f689828a12144ffa0ea8d5"
partner_vim_id = "5466c037-8ddb-47b8-b66c-a45dffe603e7"

# Initialize producer
producer = KafkaProducer(
    bootstrap_servers=bootstrap_servers,
    security_protocol=security_protocol,
    sasl_mechanism=sasl_mechanism,
    sasl_plain_username=username,
    sasl_plain_password=password,
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

# Initialize consumer
consumer = KafkaConsumer(
    RESPONSE_TOPIC,
    bootstrap_servers=bootstrap_servers,
    security_protocol=security_protocol,
    sasl_mechanism=sasl_mechanism,
    sasl_plain_username=username,
    sasl_plain_password=password,
    auto_offset_reset='latest',
    enable_auto_commit=True,
    group_id=f'perf-test-{uuid.uuid4()}',
    value_deserializer=lambda m: m.decode('utf-8'),
    consumer_timeout_ms=1000  # Poll timeout
)

# Warm up the consumer by doing an initial poll
print("Warming up consumer...")
consumer.poll(timeout_ms=2000)
print("Consumer ready\n")

def send_and_wait_response(topic, message, msg_id, timeout=120):
    """Send a message to a topic and wait for its response"""
    print(f"Sending message to {topic}...")

    # Clear any pending messages before sending
    consumer.poll(timeout_ms=100)

    producer.send(topic, message)
    producer.flush()
    print(f"✓ Message sent to {topic}\n")

    print(f"Waiting for response on {RESPONSE_TOPIC}...\n")

    start_time = time.time()
    while True:
        # Check timeout
        if time.time() - start_time > timeout:
            print(f"⚠ WARNING: Timeout after {timeout}s waiting for response")
            return None

        # Poll for messages with a short timeout
        msg_batch = consumer.poll(timeout_ms=1000)

        for topic_partition, messages in msg_batch.items():
            for response_msg in messages:
                try:
                    response = json.loads(response_msg.value)

                    if response.get('msg_id') == msg_id:
                        print(f"Received response")
                        print(f"Response: {json.dumps(response, indent=2)}\n")
                        return response
                    else:
                        print(f"Skipping message with different msg_id: {response.get('msg_id')}")

                except json.JSONDecodeError as e:
                    print(f"\n⚠ ERROR: Failed to parse message: {response_msg.value}")
                    print(f"Error details: {e}")

def onboard_app():
    """Send app onboarding request and return the response"""
    msg_id = str(uuid.uuid4())

    message = {
        "msg_id": msg_id,
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
    }

    print(f"[ONBOARD] Message ID: {msg_id}")
    response = send_and_wait_response(ONBOARD_TOPIC, message, msg_id)

    if response:
        app_id = response.get('app_id')
        if app_id:
            print(f"App ID: {app_id}\n")
        else:
            print("Warning: No app_id in response\n")

    return response

def remove_app(app_id):
    """Send app removal request and return the response"""
    msg_id = str(uuid.uuid4())
    message = {
        "msg_id": msg_id,
        "federation_context_id": federation_context_id,
        "app_id": app_id,
    }

    print(f"[REMOVE] Message ID: {msg_id}")
    response = send_and_wait_response(REMOVE_TOPIC, message, msg_id)

    return response

def instantiate_app():
    """Send app instantiation request and return the response"""
    msg_id = str(uuid.uuid4())

    # Load config from vars.yaml
    script_dir = os.path.dirname(os.path.abspath(__file__))
    vars_file = os.path.join(script_dir, 'vars.yaml')

    with open(vars_file, 'r') as f:
        config = yaml.safe_load(f)

    message = {
        "msg_id": msg_id,
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
        "vim_id": partner_vim_id,
        "config": yaml.dump(config),
    }

    print(f"[INSTANTIATE] Message ID: {msg_id}")
    response = send_and_wait_response(INSTANTIATE_TOPIC, message, msg_id)

    if response:
        app_instance_id = response.get('app_instance_id')
        if app_instance_id:
            print(f"App instance ID: {app_instance_id}\n")
        else:
            print("Warning: No app_instance_id in response\n")

    return response


def terminate_app(app_instance_id):
    """Send app termination request and return the response"""
    msg_id = str(uuid.uuid4())

    message = {
        "msg_id": msg_id,
        "federation_context_id": federation_context_id,
        "app_instance_id": app_instance_id,
    }

    print(f"[TERMINATE] Message ID: {msg_id}")
    response = send_and_wait_response(TERMINATE_TOPIC, message, msg_id)

    return response


def main():
    try:
        # Run 30 loops of the instantiate/terminate cycle
        for loop_num in range(1, 31):
            print(f"\n{'='*60}")
            print(f"LOOP {loop_num}/30")
            print(f"{'='*60}\n")

            # Onboard app
            onboard_response = onboard_app()

            # Remove app if onboarding succeeded
            if onboard_response:
                app_id = onboard_response.get('app_id')
                if app_id:
                    time.sleep(5)
                    remove_app(app_id)
                else:
                    print(f"Loop {loop_num}: Failed to get app_id, skipping removal\n")
            else:
                print(f"Loop {loop_num}: Onboarding failed, skipping removal\n")

            # # Instantiate app
            # instantiate_response = instantiate_app()

            # # Terminate app if instantiation succeeded
            # if instantiate_response:
            #     app_instance_id = instantiate_response.get('app_instance_id')
            #     if app_instance_id:
            #         time.sleep(10)
            #         terminate_app(app_instance_id)
            #     else:
            #         print(f"Loop {loop_num}: Failed to get app_instance_id, skipping termination\n")
            # else:
            #     print(f"Loop {loop_num}: Instantiation failed, skipping termination\n")

            # Brief pause between loops
            if loop_num < 30:
                time.sleep(10)

    except KeyboardInterrupt:
        print("\nStopped by user")
    except Exception as e:
        print(f"\n⚠ ERROR: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        input("\nPress Enter to continue...")
    finally:
        producer.close()
        consumer.close()


if __name__ == "__main__":
    main()
