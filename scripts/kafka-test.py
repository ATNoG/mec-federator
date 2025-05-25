from kafka import KafkaConsumer
import sys, os

# List your topics here
# topics = ['new_app_pkg', 'delete_app_pkg', 'responses', 'cluster', 'mec-apps']
topics = ['responses', 'new_app_pkg']

# Kafka configuration
bootstrap_servers = ['10.255.41.81:31999'] 
username = os.getenv('KAFKA_USERNAME', 'user1')
password = os.getenv('KAFKA_PASSWORD', 'password')
sasl_mechanism = 'PLAIN' 
security_protocol = 'SASL_PLAINTEXT' 

group_id = 'test-client'                # optional: assign a group ID

print(username, password)

def main():
    try:
        consumer = KafkaConsumer(
            *topics,
            bootstrap_servers=bootstrap_servers,
            security_protocol=security_protocol,
            sasl_mechanism=sasl_mechanism,
            sasl_plain_username=username,
            sasl_plain_password=password,
            auto_offset_reset='latest',
            enable_auto_commit=True,
            group_id=group_id,
            value_deserializer=lambda m: m.decode('utf-8')
        )

        print(f"Available topics:")
        print(consumer.topics())

        print(f"\nSubscribed to topics: {topics}")
        print("Consuming messages. Press Ctrl+C to stop.\n")

        for message in consumer:
            print(f"[{message.topic}] {message.value}")

    except KeyboardInterrupt:
        print("\nStopped by user.")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
    finally:
        if 'consumer' in locals():
            consumer.close()

if __name__ == "__main__":
    main()
