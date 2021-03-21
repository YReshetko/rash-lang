let some_constant = 144

let doubled = fn(){
	return some_constant * 2;
}

let len = fn(v){
	return eval("sys", "len", v);
}
let time = fn(){
	return eval("sys", "time");
}

let fib = fn(val){
	if (val == 1) {
		return 0;
	}
	if (val == 2) {
		return 1;
	} else {
		return fib(val - 2) + fib(val - 1);
	}
}