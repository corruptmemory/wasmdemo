# Go (Tiny) WASM demo

A small demo program written in Go that uses the [TinyGo](https://tinygo.org/) compiler.

This project is for learning purposes.  What I wanted to do was

* Get more comfortable using Go to target Web Assembly
* Understand how to interact with the DOM from Web Assembly
* Minimize the size of the WASM produced
* Use the features of Go to make Web applications easier to write

For the most part I've accomplished what I set out to do.  A few observations:

* _Tinygo for Web Assembly is viable:_ Completely competent tool chain (nice work all the Tinygo folks!).  
* Tinygo == no stack traces: When building with the regular `go` build tool chain `panic`s produce as good a stack trace as regular go projects.  Massively helpful.
* Building with regular `go` == 1+ MB (even after running `wasm-opt`) wasm file sizes.  Tinygo == ~180KB unoptimized and ~50KB optimized wasm files
* Interfacing to the DOM via github.com/Nerzal/tinydom is mostly intuitive.  Autocomplete in IDEs works (with some setup).  Some things I had to figure out via trial and error.
* The "runtime" glue in the `wasm_exec*.js` is interesting to look at.  It's also kinda remarkable how much glue is needed.  Biggish concerns is the possibility of this glue breaking with updates to the compilers
* Go "stuff" (data structures and values) != JS data structures.  Which makes sense, but it does mean that you cannot simply "decorate" DOM elements with "Go values".  Just something to keep in mind.
* It is really, really nice that Go has true first-class closures.  Sadly, there is some incidental complexity when passing Go closures, functions, and methods to the "JavaScript" world: you are responsible for "cleaning up" the mapping objects otherwise memory will get wasted.

## Cross-language comparisons

I know that Rust has a "new hotness" w.r.t. WASM (and other things) right now, but try as I might I find no joy in Rust.  Just not my jam, I guess.  That said
my _guess_ is that if you didn't blindly download all the Rust code on the internet to build a simple Web Assembly application you could make _smaller_ (in terms
of the resulting WASM code) versions of this sample than what one can do with Tinygo.  But that's a guess.

OTOH, the language that I *do fancy*, Odin, I'm _sure_ could do the same thing this application is doing with WASM files about 1/10 the size (optimized).

## Why size matters

This is a philosophical point, but it is motivated by real-life experience with using systems in general, but the Web in particular.

The "state of the art for Web development" leaves a lot to be desired.  Slow page loads, tons of requests to servers (CDNs or otherwise), bloated JS and CSS libraries
make for an experience that could only be described as "sluggish".  Mind you, I have pretty decent internet access.  I spent some time in Firefox turning on various 
throttling levels for various popular sites and the experience is, by and large: terrible.  I don't want to add my efforts to this way of experiencing the web.

*If* you're going to need code to make your web site work, then focus on 2 things

1. Deliver the _least_ amount of code your imagination can concoct.  Do regular measurements of the "cost" (to users) of using your site.  These measurements should be run using various throttling modes because not everyone will have blazing access speeds all the time.
2. Minimize the number of "fetches of content".  No, I'm not saying return to the days of yore where one encoded graphical sprites (although that technique still works) in one large image, but be aware of the round-trip overhead of your site.  And yes, reuse assets as much as possible to increase cache hit ratios.


## Building

### Prerequisites

* The `which` command
* Go (on your path) 1.17 and higher
* Tinygo (on your path) 0.22 and higher

If you want to try optimizing your build you will need to do either:

1. Install the Web Assembly tools from here: https://github.com/WebAssembly/binaryen.  If you do this then the `which` command should be able to find it (e.g.: `which wasm-opt` should return a happy result) 
2. Clone https://github.com/WebAssembly/binaryen and build the tools locally.  If you do this then you will need to supply the `-w` argument to `build.sh` with the path to where the `wasm-opt` executable can be found.

### Build script

```bash
./build.sh -h
Usage: build.sh [-h|--help] [-c|--clean] [-C|--clean-all]
                [-g|--use-go] [-b|--build] [-o|--optimize]
                [-w|--path-to-wasm-opt]

    Build wasmdemo.

Arguments:
  -h|--help                     This help text
  -c|--clean                    Clean generated artifacts.
  -C|--clean-all                Clean all the artifacts and the Go module cache.
  -g|--use-go                   Build using regular go
  -o|--optimize                 Run 'wasm-opt' to size optimize the build
  -w|--path-to-wasm-opt <path>  Path to the 'wasm-opt' program
```

### Typical build and run

```bash
$ ./build.sh && ./server
```

The build command should build the WASM program and the sample host server

Now you can hit the sample application at `http://localhost:9090/`

Enjoy!
