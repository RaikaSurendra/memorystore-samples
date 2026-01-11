# Go Caching Sample - Learning Guide

Welcome to the learning guide for the Google Cloud Memorystore Caching Sample (Go). This documentation is designed to walk you through the design, architecture, and technical intricacies of the application.

## Table of Contents

1.  [**Architecture & Design Philosophy**](./01_architecture_design.md)
    *   The MVC Pattern in Go
    *   Dependency Injection
    *   Layered Architecture (Repo -> Service -> Controller)

2.  [**Go Language Paradigms**](./02_go_paradigms.md)
    *   Structs and Interfaces
    *   Error Handling Philosophy
    *   Context Management (`context.Context`)
    *   The Gin Web Framework

3.  [**Caching, Redis, & Valkey**](./03_caching_redis_valkey.md)
    *   The Cache-Aside Pattern
    *   Redis vs. Valkey
    *   Key Eviction & TTL
    *   The Redis Protocol (RESP)

4.  [**Database & Deployment**](./04_database_docker.md)
    *   PostgreSQL with `pgx`
    *   Docker & Multi-Stage Builds
    *   Docker Compose Orchestration

---
*Created for the Memorystore Samples Project.*
