<div style="text-align:center;">
  <img src="https://dkstatics-public.digikala.com/digikala-products/117356718.jpg?x-oss-process=image/resize,m_lfit,h_800,w_800/format,webp/quality,q_90" alt="Rate Limiting">
</div>

# Rate Limiting Algorithms with Lua and Redis

This repository contains Lua implementations of three popular rate limiting algorithms: Generic Cell Rate Algorithm (GCRA), Leaky Bucket, and Sliding Window, using Redis as the data store for rate limiting.

## Algorithms

### 1\. Generic Cell Rate Algorithm (GCRA)

The Generic Cell Rate Algorithm (GCRA) is a token bucket algorithm used for traffic shaping and rate limiting. It allows bursts of traffic up to a certain size, with the rate of the traffic regulated over time.

**Pros:**

- Simple to implement and understand.
- Provides a smooth traffic shaping mechanism.
- Can handle bursty traffic efficiently.

**Cons:**

- Requires a more complex implementation compared to simpler algorithms like Leaky Bucket.
- May not be suitable for scenarios requiring strict adherence to a precise rate.

### 2\. Leaky Bucket

The Leaky Bucket algorithm is another popular rate limiting technique where requests are processed at a constant rate, preventing bursts of traffic beyond a certain threshold. If incoming requests exceed the bucket capacity, they are either delayed or discarded.

**Pros:**

- Straightforward implementation.
- Effective in controlling traffic bursts.
- Can handle varying request rates gracefully.

**Cons:**

- May introduce delay for requests during bursty traffic periods.
- Requires tuning of bucket size and leak rate for optimal performance.

### 3\. Sliding Window

The Sliding Window algorithm maintains a window of recent requests and allows only a certain number of requests within a predefined time window. It dynamically adjusts to changes in traffic patterns, allowing for more flexibility compared to static rate limiters.

**Pros:**

- Adaptable to varying traffic patterns.
- Provides more precise control over request rate.
- Suitable for dynamic environments with fluctuating traffic.

**Cons:**

- Increased complexity in implementation compared to static rate limiters.
- May require additional memory to store the sliding window.
