package example

# Unless otherwise defined, allow is false
default allow := false

allow if {
	count(violation) == 0
}

# a server is in the violation set if...
violation contains server.id if {
	some server			# There is some server...
	public_servers[server]  	# In the public_servers set and...
	server.protocols[_] == "http"	# It contains the insecure "http" protocol.
}

# a server is in the violation set if...
violation contains server.id if {
	server := input.servers[_]	# It exists in the input.servers collection and...
	server.protocols[_] == "telnet" # It contains the "telnet" protocol.	
}

# a server exists in the 'public_servers' set if...
public_servers contains server if {
	some i, j # some pair i,j verifies that...
	# it exists in the input.servers collection and...
	server := input.servers[_]	
	# it references a port in the input.ports collection and...
	server.ports[_] == input.ports[i].id 
	# the port references a network in the input.networks collection and...
	input.ports[i].network == input.networks[j].id	
	# the network is public
	input.networks[j].public
}
