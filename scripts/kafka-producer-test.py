from kafka import KafkaProducer
import sys, os
import json

# Kafka configuration
bootstrap_servers = ['10.255.41.81:31999'] 
username = os.getenv('KAFKA_USERNAME', 'user1')
password = os.getenv('KAFKA_PASSWORD', 'password')
sasl_mechanism = 'PLAIN' 
security_protocol = 'SASL_PLAINTEXT' 

# Hardcoded messages to send
federation_context_id = "dccf4231-e08f-4302-8157-0ead06bddcab"
app_pkg_id = "685f2596232d0ef0a39b7506"

messages = {
    "new_federation": {
        "msg_id": "1",
        "client_id": "operator-b",
        "client_secret": "78H0JMNA7FyyS2waNL13omQEsmWvEHyA",
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
    }
}

# Target topic to send message to
target_topic = 'federation_new_artefact'  # Change this to send to different topics

print(f"Username: {username}")
print(f"Target topic: {target_topic}")
print(f"Message to send: {json.dumps(messages[target_topic], indent=2)}")

def main():
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