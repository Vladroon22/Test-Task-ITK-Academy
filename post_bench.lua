-- bench_post.lua
local wallet_ids = {
    "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12",
    "c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13",
    "d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
    "e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15",
    "f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16"
}

-- (0.01 -- 1000.00)
function random_amount()
    return math.random(1, 100000) / 100
end

function random_operation()
    return math.random() > 0.3 and "DEPOSIT" or "WITHDRAW"
end

request = function()
    local wallet_id = wallet_ids[math.random(#wallet_ids)]
    local operation = random_operation()
    local amount = random_amount()
    
    local headers = {}
    headers["Content-Type"] = "application/json"
    
    local body = string.format('{"uuid": "%s", "type": "%s", "balance": %f}', wallet_id, operation, amount)
    
    return wrk.format("POST", "/api/v1/wallet", headers, body)
end