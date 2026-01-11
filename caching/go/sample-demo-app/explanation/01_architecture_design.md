# Chapter 1: Architecture & Design Philosophy

## Introduction
Welcome! This guide is written for developers who may be new to Go, Microservices, or the specific patterns used in cloud-native applications. We will break down *how* this application is built and *why* we made those choices.

## 1. The MVC Pattern (Model-View-Controller)

You might have heard of MVC. It is a standard way to organize code so that different parts of your application don't get tangled up "spaghetti code."

### The "Spaghetti" Problem
Imagine a single file that handles a user clicking a button, calculates a tax rate, connects to the database, sends an email, and then updates the HTML. If you want to change the email provider, you risk breaking the tax calculation.

### The MVC Solution
We split the code into three distinct roles:

1.  **Model (`/models`)**: **The Data**.
    *   This is the "shape" of your information.
    *   In our app, the `Item` struct is the model. It defines that an item has an `ID`, `Name`, `Price`, etc.
    *   It doesn't know about HTML or Databases. It's just data.

2.  **View (`/templates`)**: **The User Interface**.
    *   This is what the user *sees*.
    *   In our app, `index.html` is the view. It displays data but doesn't know how to calculate it or save it.

3.  **Controller (`/controllers`)**: **The Brains**.
    *   This accepts input (User clicks "Search").
    *   It asks the Service/Model for data.
    *   It sends that data to the View.

## 2. The Layered Architecture

In modern backend development, we often add more layers to MVC to separate "Business Logic" from "Infrastructure".

### The Flow of a Request
When you click "Search" for Item #5, here is exactly what happens:

```mermaid
graph TD
    User((User)) -->|HTTP GET /api/item/5| Router[Gin Router]
    Router -->|Calls| ItemController[ItemController]
    ItemController -->|Calls| DataController[DataController / Service]
    
    subgraph Service Layer "The Brain"
    DataController -->|1. Check Cache| Redis[(Redis/System)]
    DataController -->|2. If Miss, Check DB| Repo[ItemsRepository]
    end
    
    Repo -->|SQL Query| Postgres[(PostgreSQL)]
    
    DataController -->|Returns Item| ItemController
    ItemController -->|JSON Response| User
```

### The Layers Explained

1.  **The API Layer (Controller)**
    *   **File**: `controllers/item_controller.go`
    *   **Job**: "Front Desk". It speaks `HTTP`. It knows about Status Codes (200 OK, 404 Not Found) and JSON. It doesn't know about SQL or Cache keys.

2.  **The Service Layer (DataController)**
    *   **File**: `controllers/data_controller.go`
    *   **Job**: "Manager". It knows the *rules* of the business.
    *   **Rule**: "Always check the Cache first. If it's not there, get it from the DB, then save it to the Cache for next time."
    *   This layer is the most important part of this sample.

3.  **The Repository Layer**
    *   **File**: `repositories/items_repository.go`
    *   **Job**: "Librarian". It knows how to talk to the Database (PostgreSQL).
    *   It speaks `SQL`. It doesn't know about HTTP or Caching.
    *   *Benefit*: If we wanted to switch from PostgreSQL to MySQL, we would only change this file. The rest of the app wouldn't even know.

## 3. Dependency Injection (DI)

You will see code in `main.go` that looks like this:

```go
// 1. Create the Repository
repo := repositories.NewItemsRepository()

// 2. Create the Service, giving it the Repository
dataCtrl := controllers.NewDataController(repo)

// 3. Create the Controller, giving it the Service
itemCtrl := controllers.NewItemController(dataCtrl)
```

This is called **Dependency Injection**.

**Analogy**:
Imagine a Car (`Controller`). A Car needs an Engine (`Service`).
*   **Without DI**: The Car creates its own Engine inside itself using a welding torch. You can never change that engine.
*   **With DI**: You build an Engine separately and *install* it into the Car. If you want to test the Car, you can install a fake, plastic engine (Mock) to see if the wheels turn, without burning gas.

We use DI to make our code **Testable** and **Modular**.
