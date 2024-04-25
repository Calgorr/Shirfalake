local rate_limit_key = KEYS[1]
local now            = ARGV[1]
local burst          = ARGV[2]
local rate           = ARGV[3]
local period         = ARGV[4]

local emission_interval = period / rate
local burst_offset      = emission_interval * burst

local tat = redis.call("GET", rate_limit_key)

if not tat then
  tat = now
else
  tat = tonumber(tat)
end
tat = math.max(tat, now)

local new_tat = tat + emission_interval
local allow_at = new_tat - burst_offset
local diff = now - allow_at

local limited
local retry_in
local reset_in

local remaining = math.floor(diff / emission_interval) -- poor man's round