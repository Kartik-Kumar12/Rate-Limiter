-- rate_limiter.lua
-- A Lua script for rate limiting with a token bucket algorithm using floats

-- KEYS[1] is the token key (stores the remaining tokens)
-- KEYS[2] is the last refill timestamp key

-- ARGV[1] is the current timestamp
-- ARGV[2] is the refill rate (tokens per second)
-- ARGV[3] is the capacity (max tokens)

local tokens_key = KEYS[1]
local last_refill_key = KEYS[2]

local current_time = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])  -- tokens per second (float)
local capacity = tonumber(ARGV[3])     -- max tokens (float)

-- Get current token count and last refill time (treat tokens as float)
local tokens = tonumber(redis.call("GET", tokens_key) or capacity)
local last_refill_time = tonumber(redis.call("GET", last_refill_key) or current_time)

-- Calculate elapsed time since last refill
local elapsed_time = current_time - last_refill_time

-- Calculate tokens to add based on the elapsed time and refill rate (floating-point calculation)
local tokens_to_add = elapsed_time * refill_rate

-- Update the token count, but do not exceed the capacity (floating-point precision)
tokens = math.min(tokens + tokens_to_add, capacity)

-- Update the last refill time to the current time
last_refill_time = current_time

local ttl = 60 -- TTL for the token and timestamp keys

local tokens_to_set = tokens
-- Consume a token if available (supporting fractional tokens)
if tokens >= 1 then
    tokens_to_set = tokens - 1.0  -- Subtract 1 token (can be fractional)
    redis.log(redis.LOG_NOTICE, "Token consumed. Remaining tokens: " .. tokens_to_set)
else
    redis.log(redis.LOG_NOTICE, "No tokens left to consume")
end

-- Store the token count and last refill time as strings (even though they are floats)
redis.call("SETEX", tokens_key, ttl, tostring(tokens_to_set))
redis.call("SETEX", last_refill_key, ttl, tostring(last_refill_time))

-- Return the current tokens, elapsed time, and tokens added (all as floats)
return {tostring(tokens), tostring(elapsed_time), tostring(tokens_to_add)}
