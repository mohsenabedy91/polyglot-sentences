local typedefs = require "kong.db.schema.typedefs"

local predefined_permissions = {
    "NONE", "CREATE_USER", "READ_USER", "UPDATE_USER", "DELETE_USER",
    "CREATE_ROLE", "READ_ROLE", "UPDATE_ROLE", "DELETE_ROLE",
    "CREATE_PERMISSION", "READ_PERMISSION", "UPDATE_PERMISSION",
    "DELETE_PERMISSION", "ASSIGN_ROLES_TO_USER", "READ_USER_ROLES",
    "ASSIGN_PERMISSIONS_TO_ROLE", "READ_ROLE_PERMISSIONS"
}

local predefined_permissions_description = "Available permissions: " .. table.concat(predefined_permissions, ", ")

return {
    name = "ps-authorize",
    fields = {
        { protocols = typedefs.protocols_http },
        { consumer = typedefs.no_consumer },
        {
            config = {
                type = "record",
                fields = {
                    { authorize_url = typedefs.url({ default = "http://auth-service/authorize", required = true }) },
                    { token_header = typedefs.header_name { default = "Authorization", required = true } },
                    { required_permissions = {
                        type = "array",
                        description = "There are permissions configurations for APIs. " .. predefined_permissions_description,
                        len_min  = 1,
                        required = true,
                        elements = { type = "string", one_of = predefined_permissions },
                        default = { "NONE" },
                    } },
                }
            }
        }
    }
}
