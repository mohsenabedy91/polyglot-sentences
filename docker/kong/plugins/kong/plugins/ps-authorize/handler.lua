local http = require "resty.http"
local cjson = require "cjson"

local PsAuthorizeHandler = {
    VERSION = "1.0",
    PRIORITY = 1000,
}

local function send_error_response(status, message)
    kong.log.err(message)
    return kong.response.exit(status, { message = message })
end

local function get_token(conf)
    local access_token = kong.request.get_headers()[conf.token_header]
    if not access_token then
        return nil, "No access token found in request headers"
    end

    local token_value = string.match(access_token, "Bearer%s+(.*)")
    if not token_value then
        return nil, "Invalid token format"
    end

    return token_value, nil
end

local function request_authorization(conf, token_value)
    local httpc = http.new()
    local res, err = httpc:request_uri(conf.authorize_url, {
        method = "POST",
        ssl_verify = false,
        headers = {
            ["Authorization"] = "Bearer " .. token_value,
            ["Content-Type"] = "application/json",
            ["accept"] = "application/json"
        },
        body = cjson.encode({
            requiredPermissions = conf.required_permissions
        }),
    })

    if not res then
        return nil, "" .. (err or "unknown error")
    end

    return res, nil
end

local function handle_authorization_response(res)
    local authorize_response = cjson.decode(res.body)

    if authorize_response.error then
        return nil, "" .. authorize_response.error
    end

    if not authorize_response.data.authorized then
        return nil, "User is not authorized"
    end

    return authorize_response.data, nil
end

function PsAuthorizeHandler:access(conf)
    local token_value, err = get_token(conf)
    if err then
        return send_error_response(401, err)
    end

    local res, err = request_authorization(conf, token_value)
    if err then
        return send_error_response(500, err)
    end

    local data, err = handle_authorization_response(res)
    if err then
        return send_error_response(403, err)
    end

    kong.service.request.set_header('jti', data.jti)
    kong.service.request.set_header('exp', data.exp)
    kong.service.request.set_header('userID', data.id)
end

return PsAuthorizeHandler
