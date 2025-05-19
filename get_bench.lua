-- bench_get.lua
local wallet_ids = {
    "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12",
    "c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13",
    "d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
    "e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15",
    "f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16"
}

request = function()
    local wallet_id = wallet_ids[math.random(#wallet_ids)]
    local path = "/api/v1/wallet/" .. wallet_id
    local headers = {}
    headers["Content-Type"] = "application/json"
    return wrk.format("GET", path, headers)
end
