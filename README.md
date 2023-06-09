# bee

![GitHub](https://img.shields.io/github/license/jibstack64/bee) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/jibstack64/bee) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jibstack64/bee) ![GitHub all releases](https://img.shields.io/github/downloads/jibstack64/bee/total)

*An extremely rudimentary scripting language.*

**Basic example:**
```c++

buzz This is a comment;

#program;

name < @in <: "What's your name?: ";

@out <: "That's a cool name, " <+ name <+ "\n";

^program;

```

See more examples [here](https://github.com/jibstack64/bee/blob/master/tests).

> As a side-note: I do not suggest using this in your official projects. It is fairly unfinished and is quite spaghetti. I made it as a fun side-project; I might add some `raylib` bindings though!!

> You can customise the language to your liking. Just get `go` installed, change the constants in `bee.go` and you can goofisize the language as much as you'd like.

---

**bundle.py**

As an extra, I made a script that merges the interpreter with a given .bee script. This allows you to generate small standalone executables for your .bee scripts. Simply use [`bundle.py`](https://github.com/jibstack64/bee/blob/master/bundle.py).

---

**Concept:**

All statements end in a semi-colon, including comments, which are formatted as such: `buzz This is a comment;`.

To create a variable:

`my_variable < "My Variable!";`

Think of it as `<` is carrying the data on the right into the left.

To add variables:

`2 <+ 2;`

Of course, the value of this would be `4`. To actually store this in a variable, you would do: `my_variable < 2 <+ 2;`. You can add strings and numbers.

To subtract variables:

`5 <- 2;`

To multiply variables:

`5 <* 5;`

To divide variables:

`4 </ 2;`

To compare variables:

`2 <= 3` (would be `false`, of course)

Functions? Nah! But we have labels!
Here's an example:

```c++
buzz Create a label;
#bee;

buzz Say hello!;
@out <: "Hello!\n";

buzz Sleep for 1 second;
@sleep <: 1;

buzz This tells the program to go back to the label;
^bee;
```

To terminate a loop, or the program in general, you can use a `halt`, as shown:

```c++
x < 0;

buzz When this is true, the program will stop!;
stop < false;

#program;

@out <: "Hello!\n";

x < x <+ 1;

stop < x <= 10;

buzz ! tells the program to use the boolean as a halt;
!stop;
^program;
```

This program stops when `x` reaches 10 and so `Hello!` is said 10 times.

Conditions? Easy!

```c++

i < 0;

#start;

ii < @string <: i;
@out <: ii <+ "\n";

condition < i <= 10;
condition < condition <= false;

i < i <+ 1;

buzz Loop back to #start if condition is true;
!condition:start;

```

---

**Built-in functions:**

Typing:

`< @string <:` string conversion.

`< @bool <:` boolean conversion.

`< @num <:` number conversion.

`< @nil` returns nil. You can also simply define a variable as `nil`.

Time:

`@sleep <:` sleeps for the given number of seconds.

Destructors:

`@del <:` deletes the provided variable. Be careful when using this; it can cause your program to freeze.

Constructors:

`< @link <:` creates a very rough link to object(s). There is currently not method implemented to get a specific element of a link.

I/O:

`< @in <:` prompts the user with the provided string, and tunnels their input to the given variable.

`@out <:` pumps the given data to stdout.

Files:

`@load <:` loads a file from it's name. You can unload the loaded file by providing this function with an empty string.

`< @exists` gives true if the loaded file exists, otherwise false.

`@write <:` writes the provided data to the loaded file.

`< @read` reads the contents of the loaded file to a variable.

Runtime:

`@run <:` runs the string of code provided.

Holders:

`tmp <` this was originally used for debugging, however I decided to leave it in for funsies.

---

That's more-or-less it! I am definitely going to add the following within the coming weeks:
- [x] Easy condition statements.
- [ ] Ignoring labels. (essentially functions)
- [x] Bindings for more useful functions.
