# mec-federator 

## Implementation Overview

The MEC Federator is a comprehensive distributed system implemented in Go 1.23.6 that enables seamless federation between multiple Multi-access Edge Computing (MEC) operators. The system architecture follows a service-oriented design with clear separation of concerns, implementing core services for federation management, application orchestration, artifact handling, and infrastructure zone discovery. The implementation provides both East/West-Bound Interface (EWBI) compliance for inter-operator communication and a sophisticated event-driven architecture using Apache Kafka for reliable asynchronous messaging.

## Architecture and Technical Design

The federator employs a hybrid synchronous-asynchronous architecture that combines RESTful APIs for real-time operations with Kafka-based messaging for reliable distributed processing. The core federation establishment protocol implements a secure OAuth2-based handshake mechanism that enables operators to negotiate resource sharing agreements and exchange infrastructure information. The system abstracts underlying heterogeneous infrastructure through standardized zone and Virtual Infrastructure Manager (VIM) concepts, allowing applications to be deployed and migrated across different operator domains seamlessly.

The implementation features a comprehensive application lifecycle management framework supporting multiple artifact types including Helm charts, Terraform configurations, Ansible playbooks, and shell scripts. The system provides dynamic Kubernetes Deployment Unit (KDU) management capabilities, enabling real-time scaling and migration of applications between infrastructure nodes within and across federated environments. Cross-federation application deployment is facilitated through a sophisticated orchestrator service that manages the entire application lifecycle from artifact onboarding to retirement.

## Messaging and Communication Framework

The event-driven architecture utilizes Apache Kafka as the backbone for asynchronous communication, implementing specialized topics for federation lifecycle events (`new_federation`, `remove_federation`), application management (`federation_new_appi`, `federation_enable_kdu`, `federation_migrate_node`), and infrastructure monitoring (`infrastructure-info`). The system implements message correlation mechanisms with timeout-based response handling, enabling synchronous semantics over asynchronous infrastructure for critical operations while maintaining system resilience and scalability.

## Security and Performance Features

Security is enforced through OAuth2 client credentials flow with Keycloak integration, providing bearer token authentication for all inter-operator communications. The implementation includes comprehensive middleware pipelines for authentication, federation validation, and performance monitoring. Version 1.9 incorporates a specialized "results mode" designed for performance evaluation and testing, featuring detailed timing measurements, alert-based workflow tracking, and end-to-end operation latency monitoring to support empirical analysis of federation performance characteristics.

