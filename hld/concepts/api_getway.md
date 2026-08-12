If API gateway a single point how it handles single poit of failure ?

## What is API gateway ?
- it accepts the client request and route them to the correct backend service based on the API endpoint.

Is is this what load balancer do ?
    - just equally distribute the traffic to different instances.

1. API composition 
    - e.g MY orders
        - Mobile
            -  2 apis
        - PC
            - 4 api calls

2. Authentication 
    - auth2.0 flow
        - client get the token from auth serice
        - then gateway authenticate it before calling the BE services
    
3. Rate limiting 
    - manage burst limit
    - api throttling
    - api queue
        - thundering hear issue.

4. service discovery 
    - microservices can scale up and sclate down. it keeps track of the location (ip and port)


1. **Is this not a single point of failue ?**
- **DNS** → resolves domain to API Gateway endpoint  
- **Regional API Gateway** → entry point, auth, throttling  
- **ALB (Load Balancer)** → distributes traffic across AZs  
- **AZs** → fault isolation (high availability)  
- **App Servers** → stateless ideally  
- **DB** → usually multi-AZ (primary + replica)

```mermaid
flowchart LR
    A[Client] -->|DNS Query| B[Route53 / DNS]
    B -->|Resolved IP| C[Regional API Gateway]

    C -->|Forward Request| D[ALB]

    D --> E1[AZ-1]
    D --> E2[AZ-2]

    subgraph AZ-1
        E1 --> F1[EC2 / Service Pods]
    end

    subgraph AZ-2
        E2 --> F2[EC2 / Service Pods]
    end

    F1 --> G[(RDS / DB)]
    F2 --> G
