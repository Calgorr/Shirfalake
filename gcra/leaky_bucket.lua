local time_stamp  = tonumber(ARGV[1])
local rps = tonumber(ARGV[2])
local key = KEYS[1]

local min = time_stamp -1
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

local last = redis.call('ZRANGE', key, -1, -1)

local tat = time_stamp

if type(last) == 'table' and #last > 0 then
  for key,value in pairs(last) do
    tat = tonumber(value) + 1/rps

    break
  end
end