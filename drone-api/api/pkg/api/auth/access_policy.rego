package battlefield.authz

# default deny unless explicitly allowed
default allow = false

# A pilot can access a drone if the drone belongs to them
allow if {
    pilot := input.request.pilot
    drone := input.request.drone

    p := input.battlefield_auth.pilots[_]
    p.id == pilot
    drone == p.drones[_]
}

