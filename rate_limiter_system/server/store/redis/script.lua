-- rate_limiter.lua
-- A Lua script for rate limiting with a token bucket algorithm

-- KEYS[1] is the token key (stores the remaining tokens)
-- KEYS[2] is the last refill timestamp key

-- ARGV[1] is the current timestamp
-- ARGV[2] is the refill rate (tokens per second)
-- ARGV[3] is the capacity (max tokens)


local tokens_key = KEYS[1]
local last_refill_key = KEYS[2]

local current_time = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])  -- tokens per second
local capacity = tonumber(ARGV[3])     -- max tokens

-- Get current token count and last refill time
local tokens = tonumber(redis.call("GET", tokens_key) or capacity)
local last_refill_time = tonumber(redis.call("GET", last_refill_key) or current_time)

-- Calculate elapsed time since last refill
local elapsed_time = current_time - last_refill_time

-- Calculate tokens to add based on the elapsed time and refill rate
local tokens_to_add = elapsed_time * refill_rate

-- Update the token count, but do not exceed the capacity
tokens = math.min(tokens + tokens_to_add, capacity)

last_refill_time = current_time

local ttl = 600 

local tokens_to_set = tokens
-- Consume a token if available
if tokens >= 1 then
    tokens_to_set = tokens - 1
    redis.log(redis.LOG_NOTICE, "Token consumed. Remaining tokens: " .. tokens_to_set)
else
    redis.log(redis.LOG_NOTICE, "No tokens left to consume")
end
    
redis.call("SETEX", tokens_key,ttl, tokens_to_set)
redis.call("SETEX", last_refill_key, ttl,last_refill_time)

return {tokens, elapsed_time, tokens_to_add}

