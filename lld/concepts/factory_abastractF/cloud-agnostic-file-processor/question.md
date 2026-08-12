## Problem Statement

You are building a cloud-agnostic file processing service.

The service should support multiple cloud providers:
- AWS
- GCP

Each cloud provider offers a family of related services:
- Service Type AWS GCP
- Object Storage S3 GCS
- Message Queue SQS Pub/Sub
- Compute Trigger Lambda Cloud Functions

## Requirements
 - The application code must not depend on concrete cloud implementations.
 - All services used together must belong to the same cloud provider.
 - Adding a new cloud provider (e.g. Azure) should not require changes in business logic.
 - Cloud selection happens via configuration at startup.



 ## Approach

- Object creation is the key to this problem.
- we need to think of some creational design pattern to solve this problem.
- we need to created a family of objects based on a cloud provider : Abstract Factory Pattern
- 
 
 - Abstract Products
    - Storage
    - Queue
    - Compute

 - Concrete Products
    - AWS Storage
    - GCP Storage
    - AWS Queue
    - GCP Queue
    - AWS Compute
    - GCP Compute

 - abstract factory
    - Cloud Provider Factory -- <interface>
        - GetNewStorage
        - GetNewQueue
        - GetNewCompute

 - concrete factory : implemetation fo above interface.
    - AWS Factory
    - GCP Factory



classDiagram

    %% Abstract Factory
    class CloudFactory {
        <<interface>>
        +CreateObjectStorage() ObjectStorage
        +CreateMessageQueue() MessageQueue
        +CreateComputeTrigger() ComputeTrigger
    }

    %% Abstract Products
    class ObjectStorage {
        <<interface>>
        +Upload(file string)
    }

    class MessageQueue {
        <<interface>>
        +Publish(message string)
    }

    class ComputeTrigger {
        <<interface>>
        +Invoke(job string)
    }

    %% Concrete Factories
    class AWSCloudFactory
    class GCPCloudFactory

    CloudFactory <|.. AWSCloudFactory
    CloudFactory <|.. GCPCloudFactory

    %% AWS Products
    class S3Storage
    class SQSQueue
    class LambdaTrigger

    ObjectStorage <|.. S3Storage
    MessageQueue <|.. SQSQueue
    ComputeTrigger <|.. LambdaTrigger

    AWSCloudFactory --> S3Storage : creates
    AWSCloudFactory --> SQSQueue : creates
    AWSCloudFactory --> LambdaTrigger : creates

    %% GCP Products
    class GCSStorage
    class PubSubQueue
    class CloudFunctionTrigger

    ObjectStorage <|.. GCSStorage
    MessageQueue <|.. PubSubQueue
    ComputeTrigger <|.. CloudFunctionTrigger

    GCPCloudFactory --> GCSStorage : creates
    GCPCloudFactory --> PubSubQueue : creates
    GCPCloudFactory --> CloudFunctionTrigger : creates



    
    