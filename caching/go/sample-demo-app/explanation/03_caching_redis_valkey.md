# Chapter 3: Caching, Redis, & Valkey

## Introduction
This chapter explains the "Secret Sauce" of high-performance web applications: Caching.

## 1. What is Caching?

**Analogy**: 
Imagine you are a Librarian (`Application`).
*   A student asks for "Harry Potter".
*   **No Cache (Database only)**: You walk to the basement archive (`Database`), search through thousands of boxes, find the book, walk back up, and give it to the student. This takes 10 minutes.
*   **With Cache (Redis)**: You keep a small "Popular Books" shelf (`Cache`) right behind your desk. When the student asks, you reach behind you and grab it. This takes 2 seconds.

## 2. The "Cache-Aside" Pattern

There are many ways to cache. This application uses **Cache-Aside** (also called Lazy Loading). It is the most common and safest pattern.

### The Algorithm
The application (Librarian) is in charge. The Cache (Shelf) is passive.

1.  **Request**: User asks for `Item ID: 100`.
2.  **Check Cache**: "Do I have #100 on the shelf?"
    *   **HIT (Yes)**: Great! Return it immediately. (Fastest path)
    *   **MISS (No)**: Darn. Proceed to step 3.
3.  **Fetch from DB**: Go to the basement (PostgreSQL) and get #100.
4.  **Update Cache**: **Crucial Step**. Before returning, place a copy of #100 on the Cache Shelf.
    *   *Now, the next time someone asks, it will be a HIT.*
5.  **Return**: Give #100 to the user.

### Why "Aside"?
Because the application sits "aside" the cache and talks to both. The DB waits for the App to tell it what to do. The DB does not talk to the Cache directly.

## 3. Redis & Valkey

### What is Redis?
Redis is an **In-Memory Data Store**.
*   **In-Memory**: It stores data in RAM, not on the Hard Drive (SSD). accessing RAM is thousands of times faster than accessing Disk.
*   **Data Structure Store**: It's not just strings. It understands Lists, Sets, Maps (Hashes), and more.

### What is Valkey?
Recently, Redis changed its software license. In response, the open-source community (Linux Foundation) took the last truly open version of Redis and created **Valkey**.
*   It functions **exactly** the same as Redis for our purposes.
*   It uses the same commands, same Protocol (RESP), and same libraries.
*   This sample uses a `redis:alpine` docker image, but it could easily swap to `valkey/valkey:8` without changing a single line of Go code.

## 4. Time To Live (TTL)

We don't keep data in the cache forever. In `cache.go`, you see:
`Rdb.Set(ctx, key, value, 1*time.Hour)`

**Why delete data?**
1.  **Stale Data**: If the price of "Harry Potter" changes in the Database, the Cache still shows the old price. By deleting it after 1 hour, we force the App to re-fetch the new price eventually.
2.  **Memory is Expensive**: RAM is expensive and limited. We can't store the whole library on the shelf. We only want *active* items there.

## 5. Serialization

Redis only understands bytes (text/binary). It doesn't understand "Go Objects".
*   **Writing**: We must convert our Go `Item` struct into a JSON string (`json.Marshal`).
    `Go Struct -> JSON String -> Redis`
*   **Reading**: We read the string and convert it back (`json.Unmarshal`).
    `Redis -> JSON String -> Go Struct`
