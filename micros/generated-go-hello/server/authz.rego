package example.authz

default allow = false

allow { # A user is allowed if ...
    # For a defined JWT audience...
    user := input.jwt.aud
    # With a ACL...
    acl := input.access_control.users[user].acl
    acl[0] == input.uri
}
