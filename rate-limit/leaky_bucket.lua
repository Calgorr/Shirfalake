local key = KEYS[1]
local ts = tonumber(ARGV[1])
local cps = tonumber(ARGV[2])

local emission_interval = 1 / cps

local last_ts = redis.call("GET",key)
if not last_ts then
    last_ts = ts
else
    last_ts = tonumber(last_ts)
end

local limited=1
local remaining = last_ts+emission_interval-ts

if ts==last_ts then
    limited=0
    remaining=0
    redis.call("SET", key, ts, "PX", math.ceil(emission_interval))

end



return {
    limited, -- allowed
    remaining , -- remaining
}