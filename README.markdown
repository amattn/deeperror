deeperror
=========

Verbose, but very informative, enlightening and pleasantly formatted errors for Go

Installation
------------

Stop me if you've seen this before:

	go get github.com/amattn/deeperror


Basic Usage
-----------

	package main

	import (
		"github.com/amattn/deeperror"
		"log"
		"strconv"
	)

	func main() {
		innerFunc()
	}

	func innerFunc() {
		innerInnerFunc()
	}

	func innerInnerFunc() {
		_, err := strconv.Atoi("not a number!")

		derr := deeperror.New(1506851885, "Oops, we can't understand that number.  Please try again.", err)
		log.Print(derr)
	}


Sample Output
-------------

	2013/09/19 18:16:57 

	-- DeepError 1506851885 500 main.go main.innerInnerFunc line: 20 
	-- EndUserMsg:  Oops, we can't understand that number.  Please try again. 
	-- DebugMsg:   
	-- StackTrace: -- goroutine 1 [running]:
	-- github.com/amattn/deeperror.New(0x59d0bc2d, 0x2534d0, 0x39, 0x2104418a0, 0x210441870, ...)
	-- 	/Users/kai/Dropbox/gitStore/github.com/test/src/github.com/amattn/deeperror/error.go:55 +0x1cc
	-- main.innerInnerFunc()
	-- 	/Users/kai/Dropbox/gitStore/github.com/test/src/testdeeperror/main.go:20 +0x6a
	-- main.innerFunc()
	-- 	/Users/kai/Dropbox/gitStore/github.com/test/src/testdeeperror/main.go:14 +0x18
	-- main.main()
	-- 	/Users/kai/Dropbox/gitStore/github.com/test/src/testdeeperror/main.go:10 +0x18
	--  
	-- ParentError: -- strconv.ParseInt: parsing "not a number!": invalid syntax

But Why?
--------


Because debugging is hard.  I wanted to make debuggin easier.  These are the tricks I use:

1. Error Numbers
2. End User Error Messages
3. Debug Messages
4. Stack Traces
5. Error Chaining

Since you asked, 1. is my favorite, I use it in every language and platform I work with and 5. is especially powerful and relevent in Go.  What?  You didn't ask?  Oh.  Nevermind.  I may have misheard.

### Error Numbers

The err number (1506851885) makes it easy to do a find.  Line numbers are insufficent when people are making changes to code and you don't know which exact version of the app is generating which stack trace.
I use this shell command to generate them:

	#!/bin/sh
	od -vAn -N4 -tu4 < /dev/urandom | tr -d " \n"

I use Keyboard Maestro to run the command and paste in where my cursor is.  You could probably use any sufficently smart macro tool.

### End User Error Messages

It's always nice to have the option to show something to the user when things go wrong.  Sometimes it is their fault.  If only users wouldn't cause bugs, I wouldn't need this item.  The issue is that debug messages and end user messages should be different.  You log the debug messages and you only show the user the EndUserMsg.  Maybe include the err number so that support calls can be streamlined.

### Debug Messages

Because debugging is hard, it's nice to get some extra hints now and then.  I stuff this string with ancilliary data which might help me figure out what is going wrong.  These don't get presented to the user, but they do get logged.

### Stack Traces

By themselves, sometimes useful.  combined with the next point, they become even more powerful than you could ever imagine.

### Error Chaining

Go usually doesn't do exceptions.  The normal pattern is to "pass the error up" the call stack.  By chaining the errors, we can pinpoint the error source faster and usualy get a successful repro sooner.  