```
@@@@@@@    @@@@@@    @@@@@@   @@@  @@@  @@@@@@@@   @@@@@@   @@@  @@@   @@@@@@   
@@@@@@@@  @@@@@@@@  @@@@@@@   @@@  @@@  @@@@@@@@  @@@@@@@   @@@  @@@  @@@@@@@@  
@@!  @@@  @@!  @@@  !@@       @@!  @@@  @@!       !@@       @@!  !@@  @@!  @@@  
!@!  @!@  !@!  @!@  !@!       !@!  @!@  !@!       !@!       !@!  @!!  !@!  @!@  
@!@!!@!   @!@!@!@!  !!@@!!    @!@!@!@!  @!!!:!    !!@@!!    @!@@!@!   @!@!@!@!  
!!@!@!    !!!@!!!!   !!@!!!   !!!@!!!!  !!!!!:     !!@!!!   !!@!!!    !!!@!!!!  
!!: :!!   !!:  !!!       !:!  !!:  !!!  !!:            !:!  !!: :!!   !!:  !!!  
:!:  !:!  :!:  !:!      !:!   :!:  !:!  :!:           !:!   :!:  !:!  :!:  !:!  
::   :::  ::   :::  :::: ::   ::   :::   :: ::::  :::: ::    ::  :::  ::   :::  
 :   : :   :   : :  :: : :     :   : :  : :: ::   :: : :     :   :::   :   : :  
```

Is an opensource scripting programming language with an interpreter written on golang. 
The language functionality is simply extended by writing plugins.

#Language types
* Integer, defined by literal, for example ``` let a = 42;```
* Boolean, defined by next literals: `true`, `false`. For example ``` let bool = true;```
* String, defined by double quotes: `"Hello world!"`
* Array, defined by `[` from left side and `]` - from right. Values are separated by a comma, for example: ```let arr = [1, "hello", true, "world"];```
* Hash, the pairs of hashable literals separated by a comma. Each pair separated by a colon. For example: ```let map = {"one": 1, 2 : "two", true: "three"};```
* Function, defined by `fn` literal, contains a block of arguments and block of statements: ```fn(<arguments>){<statements>};``` 
For example:
  ```
  let func = fn(a, b){
    return a + b;
  }
  ```

#Statements
* `let` - creates a new variable in execution scope and assigns a value, for example: ```let a = 10;```
* Assign - assigns value to existing variable or map/array elements in scope, for example: ```a = 12; map["one"] = true; arr[10] = 50;```
* `return` - returns value from functional call. Can be omitted, because the language returns value of last execution in a block. For example: ```return 10;``` and ```10;``` are equal. The difference is that explicit `return` call can break function execution.
* Declaration - rash supports import one script files to another. The declaration starts from `#` then alias and string literal with path to the script. For example: ```# sys "lib/sys.rs"```. Then variables of imported script available by alias, for example: ```let a = sys.tick;```

#Builtin funcctions
* `eval` - allows to make a single call to plugin, signature: ```eval(<package name>, <function name>, <any number of arguments>);```
    * `package name` - string literal, defined in plugin and registered in an interpreter with that name;
    * `function name` - string literal, helps plugin to understand which function has to be called;
    * `any number of arguments` - arguments which has to be sent to particular function in the plugin;
* `call` - allows to register callback function in a plugin, signature: ```eval(<package name>, <function name>, <callback function>, <any number of arguments>);```
    * `package name` - string literal, defined in plugin and registered in an interpreter with that name;
    * `function name` - string literal, helps plugin to understand which function has to be called;
    * `callback function` - the function defined in rash language which arguments number and returned value corresponds to plugin specification
    * `any number of arguments` - arguments which has to be sent to particular function in the plugin;

#Operations
* `+` - supported on strings and integers
* `-` - supported on integers
* `*` - ...
* `/` - ...
* `==` - ...
* `!=` - ...
* `>` - ...
* `<` - ...
* `if` - classic if operator which supports two types: ```if (<condition>) {<block statements>}``` and ```if (<condition>) {<block statements>} else {<block statements>}```. Can be used as ternary operator: ```let a = if (b == c) {true} else {false}``` 

# How to extend lib
* If you don't need some special functionality you can just create a separate script file, for example `some_file.rs` where you put your code and then include it where you need: `# some "some_file.rs"`
* If you need complicated functionality the you can implement interface:
```go
type Plugin interface {
	Eval(fnName string, args ...interface{}) ([]interface{}, error)
	Call(fnName string, callback func(args ...interface{}) ([]interface{}, error), args ...interface{}) ([]interface{}, error)
	Package() string
	Version() string
	Description() string
}
```
And inject it into interpreter by modifying main.go. Also to include the functionality to your code it's better to create *.rs wrappers for each plugin, so you can naturally use the functionality in your scripts.

# Run
At the moment there is supported only one way, it is run REPL app: `make run` . You need go installed on your machine. 
