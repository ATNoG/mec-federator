from kafka import KafkaProducer
import sys, os
import json
import argparse

# Kafka configuration
bootstrap_servers = ['10.255.41.81:31999'] 
username = os.getenv('KAFKA_USERNAME', 'user1')
password = os.getenv('KAFKA_PASSWORD', 'password')
sasl_mechanism = 'PLAIN' 
security_protocol = 'SASL_PLAINTEXT' 

# Hardcoded messages to send
federation_context_id = "92853657-76ec-4dd7-860b-829f6047e520"
app_pkg_id = "686403a4232d0ef0a39b7508"

messages = {
    "new_federation": {
        "msg_id": "1",
        "client_id": "operator-a",
        "client_secret": "FkEZE8twp5sMn3qcVqvm3nZKzy9sLAr8",
        "federation_endpoint": "http://federator-po:8000",
        "auth_endpoint": "http://federator-po:8000/federation/v1/auth/token",
    },
    "remove_federation": {
        "msg_id": "2",
        "federation_context_id": federation_context_id,
    },
    "federation_new_artefact": {
        "msg_id": "3",
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
    },
    "federation_remove_artefact": {
        "msg_id": "4",
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
    }
}

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Send test messages to Kafka topics')
    parser.add_argument('topic', help='Topic name to send message to', 
                       choices=list(messages.keys()))
    
    args = parser.parse_args()
    target_topic = args.topic

    print(f"Username: {username}")
    print(f"Target topic: {target_topic}")
    print(f"Message to send: {json.dumps(messages[target_topic], indent=2)}")
    try:
        producer = KafkaProducer(
            bootstrap_servers=bootstrap_servers,
            security_protocol=security_protocol,
            sasl_mechanism=sasl_mechanism,
            sasl_plain_username=username,
            sasl_plain_password=password,
            value_serializer=lambda v: json.dumps(v).encode('utf-8')
        )

        print(f"\nAvailable topics:")
        # Note: Producer doesn't have direct access to topics like consumer
        # You can list them manually or get from metadata if needed

        print(f"\nSending message to topic: {target_topic}")
        
        # Send the message
        future = producer.send(target_topic, messages[target_topic])
        
        # Wait for the message to be sent
        record_metadata = future.get(timeout=10)
        
        print(f"Message sent successfully!")
        print(f"Topic: {record_metadata.topic}")
        print(f"Partition: {record_metadata.partition}")
        print(f"Offset: {record_metadata.offset}")

    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
    finally:
        if 'producer' in locals():
            producer.close()

if __name__ == "__main__":
    main() 
