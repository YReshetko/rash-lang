let new_server = fn(port){
	let server_name = eval("http", "new", port);
	return {
		"register": fn(method, route, function){
			call("http", "register", function, server_name, method, route);
		},
		"start": fn(){
			eval("http", "start", server_name);
		}
	}
}