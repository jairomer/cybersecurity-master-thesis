package battlefield.authz

# default deny unless explicitly allowed
default allow = false

# A drone can only have access to its own data.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role

    p := input.battlefield.pilots[_]

    role == "drone"
    p.drones[_] == user
    drone := user
}

# A pilot can access a drone if the drone belongs to them.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role

    role == "pilot"
    p := input.battlefield.pilots[_]
    p.id == user
    drone := p.drones[_]
}

# An officer can access data from all drones in the battlefield.
drones[drone] if {
    user := input.request.user.id
    role := input.request.user.role

    role == "officer"
    p = input.battlefield.pilots[_]
    drone := p.drones[_]
}

allow if {
    count(drones) > 0
}