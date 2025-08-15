trivial:
	opa eval "1+2+3"

cmd:
	opa eval -i input.json -d example.rego "data.example.violation[x]"

cmdexit:
	opa eval --fail-defined -i input.json -d example.rego "data.example.violation[x]";
