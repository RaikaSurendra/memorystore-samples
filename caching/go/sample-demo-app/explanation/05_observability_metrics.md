# Chapter 5: Observability & Metrics

## 1. Understanding Redis Memory

You asked: **"What does Redis Memory Used mean?"**

Redis metrics can be confusing. Here is the breakdown:

### `used_memory` (The Content)
*   **Definition**: The total number of bytes allocated by Redis to store your data (keys, values, internal structures).
*   **Analogy**: This is the *weight of the clothes* inside your suitcase.
*   **Why it matters**: If this hits your `maxmemory` limit, Redis will start deleting (evicting) keys to make room.

### `used_memory_rss` (The Suitcase)
*   **Definition**: Resident Set Size. The amount of memory the Operating System has allocated to the Redis process.
*   **Analogy**: This is the *size of the suitcase* itself allocated by the OS.
*   **Fragmentation**: Sometimes RSS > Used. This is like a half-empty suitcase. The space is "taken" but not fully used.

## 2. What Should Be Added to the Dashboard?

The current dashboard is a "Health" dashboard (Is it alive?). To make it a "Performance" dashboard, I recommend adding these panels:

### A. Cache Hit Ratio (The Golden Metric)
*   **Formula**: `Hits / (Hits + Misses)`
*   **Target**: > 80% is often the goal.
*   **Why**: If your hit ratio is 5%, your cache is useless. You are paying for Redis but still hitting the DB for almost every request.
*   *Required Change*: We need to add a custom Prometheus Counter in the Go code (`cache_hits_total`, `cache_misses_total`) inside `DataController`.

### B. Eviction Rate
*   **Metric**: `rate(redis_evicted_keys_total[1m])`
*   **Why**: If this number spikes, your Cache is full. Redis is frantically deleting old data to fit new data. You might need a bigger Redis instance or a shorter TTL.

### C. Database Connection Pool Saturation
*   **Metric**: `pgxpool` stats (if exposed).
*   **Why**: If your app is waiting 200ms just to *get* a connection to Postgres, your database performance doesn't matter. You are bottlenecked by the pool size.

### D. 99th Percentile Latency (P99)
*   **Metric**: `histogram_quantile(0.99, ...)`
*   **Why**: You currently have P95. P99 tells you the experience of the "unluckiest" 1% of your users. These are often the requests that failed or timed out.
