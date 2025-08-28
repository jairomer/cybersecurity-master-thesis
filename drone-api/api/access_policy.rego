package battlefield.authz

# default deny unless explicitly allowed
default allow = false

# A drone can only have access to its own data.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role
    op := input.request.user.operation

    p := input.battlefield.pilots[_]

    role == "drone"
    p.drones[_] == user
    # Allowed operations
    op == {"get-target", "set-location", "get-battlefield"}[_]
    drone := user
}

# A pilot can access a drone if the drone belongs to them.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role
    op := input.request.user.operation

    role == "pilot"
    p := input.battlefield.pilots[_]
    p.id == user
    # Allowed operations
    op == {"set-target", "get-battlefield"}[_]
    drone := p.drones[_]
}

# An officer can access data from all drones in the battlefield.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role
    op := input.request.user.operation

    role == "officer"
    p = input.battlefield.pilots[_]
    # Allowed operations
    op == {"get-battlefield", "provisioning"}[_]
    drone := p.drones[_]
}

allow if {
    count(drones) > 0
}