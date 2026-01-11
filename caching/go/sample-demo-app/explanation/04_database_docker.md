# Chapter 4: Database & Deployment

## Introduction
This chapter explains how we save data permanently and how we package the application so it runs anywhere.

## 1. PostgreSQL & Connection Pooling

We use **PostgreSQL** as our "Source of Truth". If the power goes out, Redis loses data (it's in RAM), but Postgres keeps it (it's on Disk).

### The Connection Pool (`pgxpool`)
Connecting to a database is "expensive". It involves a TCP handshake, authentication, and process creation. It takes time (milliseconds, which add up).

**Bad Application**: Opens a connection, runs one query, closes connection.
*   User 1 waits 50ms for connection + 5ms for query.
*   User 2 waits 50ms for connection + 5ms for query.

**Our Application**: Uses a **Pool**.
*   The App starts and opens 10 connections immediately. They sit IDLE, waiting.
*   User 1 asks for data. The Pool gives them connection #3. Query takes 5ms. Connection #3 goes back to the pool.
*   **Result**: Zero setup time for requests. Massive performance gain.

## 2. Docker: "Works on my Machine"

Have you ever had code that works on your laptop but crashes on the server? Docker fixes that.

### What is a Container?
Think of a **Container** as a "Lightweight Virtual Machine".
*   It contains the Operating System (Linux), the libraries, the files, and your app.
*   It is **Isolated**. It doesn't care if your laptop is Mac, Windows, or Linux. Inside the container, it's always Alpine Linux.

### The Dockerfile Explained

Our `Dockerfile` uses a **Multi-Stage Build**. This is an advanced technique for tiny images.

#### Stage 1: The Factory (Builder)
```dockerfile
FROM golang:1.24 as builder
```
*   **Context**: This image is HUGE (800MB+). It has the Go Compiler, source code tools, etc.
*   **Action**: We put our code in, run `go build`, and get a single binary file: `server`.

#### Stage 2: The Shipping Box (Runtime)
```dockerfile
FROM alpine:3.18
```
*   **Context**: This image is TINY (5MB). It has almost nothing. No Go compiler, no source code.
*   **Action**: We strictly `COPY` the `server` binary from the Factory to here.
*   **Result**: Our final "Shipping Box" is small and secure. Hackers can't read our source code because it's not even there!

## 3. Orchestration with Docker Compose

We have 3 containers:
1.  **App** (Go)
2.  **Cache** (Redis)
3.  **Database** (Postgres)

We could start them manually (`docker run ...`), but we'd have to manage their IP addresses so they can talk to each other.

**Docker Compose** is the conductor:
*   It reads `docker-compose.yml`.
*   It creates a private **Network** just for these 3.
*   It gives them names (`redis`, `postgres`).
*   It ensures the DB starts before the App.

This is why in `main.go` checking, we can say `os.Getenv("REDIS_HOST")` is "redis", and it magically finds the right IP address.
