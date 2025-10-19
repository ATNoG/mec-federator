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
# federation_endpoint = "http://federator-po:8000"
federation_endpoint = "http://10.255.41.64:32440"
auth_endpoint = federation_endpoint + "/federation/v1/auth/token"

federation_context_id = "cbeb895e-64c9-496e-a77f-670728d1fcc0"

zone_id = "a61c4ca9-c7ce-4cd4-a3c8-98511fa75556"

artefact_id = "55314322-34a6-405e-81ea-1595c92f8b80"

app_id = "c0e5e719-1179-4fdd-908f-ea94e6bf0a01"

app_instance_id = "16c114dd-aea3-4ac0-a588-bad4e100b01a"

app_pkg_id = "68f12328eae7e7de1b70711b"
partner_vim_id = "45af887d-7fef-4c82-9428-d75fe43108e8"
mec_appd_id = "mec-test-server-appd"
ns_id = "203972c9-ccac-4a14-8df3-63a7ee782e5b"
vnf_id = "283b1534-89ab-42f0-9fb1-cf2284010801"
kdu_id = "mec-test-server"

messages = {
    "new_federation": {
        "msg_id": "1",
        "client_id": "operator-a",
        "client_secret": "FkEZE8twp5sMn3qcVqvm3nZKzy9sLAr8",
        "federation_endpoint": federation_endpoint,
        "auth_endpoint": auth_endpoint,
    },
    "remove_federation": {
        "msg_id": "2",
        "federation_context_id": federation_context_id,
    },
    "federation_get_info":{
        "msg_id": "22",
        "federation_context_id": federation_context_id,
    },
    "federation_subscribe_zone": {
        "msg_id": "23",
        "federation_context_id": federation_context_id,
        "zone_id": zone_id,
    },
    "federation_unsubscribe_zone": {
        "msg_id": "24",
        "federation_context_id": federation_context_id,
        "zone_id": zone_id,
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
    },
    "federation_new_app": {
        "msg_id": "5",
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
    },
    "federation_remove_app": {
        "msg_id": "6",
        "federation_context_id": federation_context_id,
        "app_id": app_id,
    },
    "federation_new_appi": {
        "msg_id": "5",
        "federation_context_id": federation_context_id,
        "app_pkg_id": app_pkg_id,
        "vim_id": partner_vim_id,
        "config": "",
    },
    "federation_get_appi_info": {
        "msg_id": "6",
        "federation_context_id": federation_context_id,
        "app_instance_id": app_instance_id,
    },
    "federation_remove_appi": {
        "msg_id": "6",
        "federation_context_id": federation_context_id,
        "app_instance_id": app_instance_id,
    },
    "federation_enable_kdu": {
        "msg_id": "7",
        "federation_context_id": federation_context_id,
        "mec_appd_id": mec_appd_id,
        "ns_id": ns_id,
        "kdu_id": kdu_id,
        "node": "cluster2-96720290",
    },
    "federation_disable_kdu": {
        "msg_id": "8",
        "federation_context_id": federation_context_id,
        "mec_appd_id": mec_appd_id,
        "ns_id": ns_id,
        "kdu_id": kdu_id,
    },
    "federation-infrastructure-info": {
        "msg_id": "9",
        "federation-meh-metrics": {
            "cluster1Id": "cluster1",
            "cluster2Id": "cluster2",
            "cluster3Id": "cluster3",
        },
    },
    "federation_migrate_node": {
        "msg_id": "10",
        "federation_context_id": federation_context_id,
        "ns_id": ns_id,
        "vnf_id": vnf_id,
        "kdu_id": kdu_id,
        "node": "cluster2-96720290",
    },
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
