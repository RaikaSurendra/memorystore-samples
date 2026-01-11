# Chapter 2: Go Language Paradigms

## Introduction
Go (or Golang) is famous for being simple, fast, and opinionated. If you are coming from Java, Python, or JavaScript, some things might look weird. This chapter explains the "Go way" of doing things found in this project.

## 1. Structs vs. Classes

Go is not an Object-Oriented language in the traditional sense. It does not have `class`, `extends`, or `inheritance`.

### The Struct
Instead, we use `struct`. A struct is just a collection of fields.

```go
type Item struct {
    ID          int64
    Name        string
    Description string
}
```

### Methods
We can attach functions to a struct. These are called **Methods**.

```go
// The (i *Item) part is the "Receiver". It says "this function belongs to Item".
func (i *Item) ToJSON() string {
    // ...
}
```
**Analogy**: Think of a `struct` as a custom data form, and `methods` as the specific tools designed to work on that form.

## 2. Interfaces (The "Can Do" Pattern)

In Java, you explicitly say `class MyRepo implements Repository`.
In Go, interfaces are **implicit**. You don't "implement" them. You just *do* them.

If an interface asks for:
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
And your unrelated struct happens to have a `Read` method with that exact signature, Go automatically treats it as a `Reader`.
*   *Why?* It allows for very loose coupling. You can mock anything if you match the method signature.

## 3. Error Handling (No Exceptions!)

This is the most controversial and distinctive part of Go.
**Go does not have Exceptions.** No `try-catch` blocks that jump across the code.

### The Pattern
Functions return **two** things: the Result and an Error.

```go
result, err := database.Query("SELECT * ...")
if err != nil {
    // SOMETHING BAD HAPPENED!
    // We must handle it right here, right now.
    log.Println("DB failed:", err)
    return
}
// If we are here, we know it worked.
use(result)
```

**Philosophy**: Exceptions are "magical" flow control (GOTO statements). They make it hard to see where the program might crash. Go forces you to acknowledge that **failure is a normal part of detailed software operations**, not an "exception".

## 4. Context (`context.Context`)

You will see `ctx` or `context` passed as the first argument to almost every function that talks to an external system (DB, Redis).

**The Problem**:
User A requests a heavy report. The DB starts working. User A gets bored and closes the browser tab.
*   **Without Context**: The DB keeps working for 5 minutes, wasting resources, then realizes nobody is listening.
*   **With Context**: The web server notices the connection closed. It cancels the `Context`. The `Context` signal travels down to the DB driver, which immediately kills the query.

**Analogy**:
Context is like a "Mission Time Limit" combined with a "Walkie Talkie".
*   "You have 5 seconds to get the data" (Timeout).
*   "Abort mission! The user left!" (Cancellation).

## 5. JSON Tags

```go
type Item struct {
    Name string `json:"name" binding:"required"`
}
```

The text between backticks `` `...` `` is a **Tag**.
*   `json:"name"`: Tells the JSON encoder: "When you send this to the browser, rename `Name` (uppercase) to `name` (lowercase)".
*   `binding:"required"`: Tells the Gin framework: "If the user sends a POST request without a name, reject it immediately with an error".
