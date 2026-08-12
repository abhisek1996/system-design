# OOP Concepts – Exhaustive Interview Notes

---

## 1. The 4 Pillars of OOP

### 1.1 Encapsulation
Bundling data (fields) and methods that operate on that data into a single unit (class), and **restricting direct access** to internal state.

- Achieved using **access modifiers**: `private`, `protected`, `public`
- Data is accessed via **getters and setters**
- Protects object integrity by preventing external code from putting it in an invalid state

```java
class BankAccount {
    private double balance; // hidden from outside

    public double getBalance() { return balance; }

    public void deposit(double amount) {
        if (amount > 0) balance += amount; // validation inside
    }
}
```

**Interview tip:** Encapsulation is about *protection and control*, not just grouping. Emphasize validation logic in setters.

---

### 1.2 Abstraction
Hiding **implementation complexity** and exposing only what's necessary to the user.

- Achieved via **abstract classes** and **interfaces**
- The user knows *what* an object does, not *how* it does it

```java
abstract class Shape {
    abstract double area(); // user just calls area(), doesn't care about impl
}

class Circle extends Shape {
    double radius;
    double area() { return Math.PI * radius * radius; }
}
```

**Encapsulation vs Abstraction:**
| | Encapsulation | Abstraction |
|---|---|---|
| Focus | Hiding *data* | Hiding *implementation* |
| Achieved by | Access modifiers | Abstract classes / Interfaces |
| Purpose | Protect internal state | Reduce complexity for the user |

---

### 1.3 Inheritance
A class (**child/subclass**) acquires properties and behaviors of another class (**parent/superclass**).

- Promotes **code reuse**
- Represents an **"is-a"** relationship
- `extends` keyword in Java/Kotlin; `:` in C#/C++

```java
class Animal {
    void eat() { System.out.println("Eating..."); }
}

class Dog extends Animal {
    void bark() { System.out.println("Barking..."); }
}

Dog d = new Dog();
d.eat();  // inherited
d.bark(); // own method
```

**Types of Inheritance:**
- **Single** – one parent, one child
- **Multilevel** – A → B → C
- **Hierarchical** – one parent, multiple children
- **Multiple** – multiple parents (not supported in Java for classes; supported via interfaces)
- **Hybrid** – combination of above

**Why Java doesn't support multiple inheritance with classes?**
The **Diamond Problem** — if two parent classes have the same method, the child doesn't know which to inherit. Java solves this by allowing multiple **interface** implementation instead.

---

### 1.4 Polymorphism
The ability of an object to take **many forms** — the same interface behaves differently based on the actual object type.

**Two types:**

#### Compile-time (Static) Polymorphism — Method Overloading
Same method name, different parameters. Resolved at compile time.

```java
class Calculator {
    int add(int a, int b) { return a + b; }
    double add(double a, double b) { return a + b; } // overloaded
}
```

#### Runtime (Dynamic) Polymorphism — Method Overriding
Child class provides its own implementation of a parent method. Resolved at runtime.

```java
class Animal {
    void sound() { System.out.println("Some sound"); }
}
class Cat extends Animal {
    @Override
    void sound() { System.out.println("Meow"); }
}

Animal a = new Cat(); // parent reference, child object
a.sound(); // prints "Meow" — runtime decision
```

**Interview tip:** The key to runtime polymorphism is **parent reference pointing to child object**. This is the basis of the Open/Closed principle (SOLID).

---

## 2. Classes and Objects

- **Class** – a blueprint/template
- **Object** – an instance of a class, created at runtime in heap memory
- **Constructor** – special method called when object is created; same name as class, no return type

```java
class Car {
    String brand;

    Car(String brand) {      // constructor
        this.brand = brand;
    }
}

Car c = new Car("Toyota"); // object creation
```

**Types of constructors:**
- **Default constructor** – no args; auto-provided if none defined
- **Parameterized constructor** – accepts arguments
- **Copy constructor** – creates a new object from an existing one (explicit in Java)

---

## 3. Abstract Classes vs Interfaces

| Feature | Abstract Class | Interface |
|---|---|---|
| Methods | Can have abstract + concrete methods | Only abstract (Java 7); default/static allowed from Java 8 |
| Variables | Can have instance variables | Only `public static final` constants |
| Constructor | Can have constructors | Cannot |
| Inheritance | Single inheritance only | Can implement multiple |
| When to use | Shared base behavior (partial impl) | Define a contract / capability |

```java
// Abstract class
abstract class Vehicle {
    String brand;
    abstract void startEngine();
    void stop() { System.out.println("Stopping"); } // concrete method
}

// Interface
interface Flyable {
    void fly(); // abstract by default
    default void land() { System.out.println("Landing"); } // default method (Java 8+)
}
```

**Interview tip:** "Use abstract class when classes share common code. Use interface when unrelated classes need to implement the same capability (e.g., `Serializable`, `Comparable`)."

---

## 4. Method Overloading vs Overriding

| | Overloading | Overriding |
|---|---|---|
| Type | Compile-time polymorphism | Runtime polymorphism |
| Class | Same class | Parent-child classes |
| Signature | Must differ | Must be same |
| Return type | Can differ | Must be same (or covariant) |
| Access modifier | Any | Cannot be more restrictive |
| `static` methods | Can be overloaded | Cannot be overridden (hidden instead) |

---

## 5. Access Modifiers

| Modifier | Same Class | Same Package | Subclass | Outside |
|---|---|---|---|---|
| `private` | ✅ | ❌ | ❌ | ❌ |
| `default` (no keyword) | ✅ | ✅ | ❌ | ❌ |
| `protected` | ✅ | ✅ | ✅ | ❌ |
| `public` | ✅ | ✅ | ✅ | ✅ |

---

## 6. `this` and `super` Keywords

### `this`
- Refers to the **current object**
- Used to differentiate instance variables from local variables
- Can call another constructor in the same class: `this(args)`

### `super`
- Refers to the **parent class**
- Used to call parent class constructor: `super(args)` — must be first line
- Used to call overridden parent method: `super.methodName()`

```java
class Animal {
    String name;
    Animal(String name) { this.name = name; }
    void speak() { System.out.println("..."); }
}

class Dog extends Animal {
    Dog(String name) {
        super(name); // calls Animal constructor
    }
    void speak() {
        super.speak(); // calls Animal's speak
        System.out.println("Woof");
    }
}
```

---

## 7. `static` Keyword

- **Static variable** – shared across all instances (class-level, not object-level)
- **Static method** – can be called without creating an object; cannot access instance variables
- **Static block** – runs once when class is loaded; used for static initialization
- **Static nested class** – doesn't need an outer class instance

```java
class Counter {
    static int count = 0; // shared across all objects

    Counter() { count++; }

    static int getCount() { return count; } // called as Counter.getCount()
}
```

**Interview tip:** `static` methods cannot be overridden — they are *hidden*, not overridden. This is because polymorphism works through objects, and static methods belong to the class.

---

## 8. `final` Keyword

| Usage | Meaning |
|---|---|
| `final` variable | Value cannot be changed (constant) |
| `final` method | Cannot be overridden |
| `final` class | Cannot be subclassed (e.g., `String` in Java) |

---

## 9. Constructors – Deep Dive

- Not inherited by subclasses
- Cannot be `static`, `abstract`, or `final`
- If a parent class has a parameterized constructor and no default constructor, the child **must** explicitly call `super(args)`
- **Constructor chaining** – calling one constructor from another using `this()` or `super()`

---

## 10. Object Class Methods (Java)

Every class implicitly extends `Object`. Key methods to know:

| Method | Purpose |
|---|---|
| `toString()` | String representation of object |
| `equals()` | Logical equality (override for value comparison) |
| `hashCode()` | Used in hash-based collections; must be consistent with `equals()` |
| `clone()` | Creates a copy (shallow by default) |
| `getClass()` | Returns runtime class |

**Important rule:** If you override `equals()`, you **must** override `hashCode()` too. Two objects that are `equal` must have the same `hashCode`.

---

## 11. Composition vs Inheritance

| | Inheritance | Composition |
|---|---|---|
| Relationship | "is-a" | "has-a" |
| Coupling | Tight | Loose |
| Flexibility | Less flexible | More flexible |
| Example | `Dog is an Animal` | `Car has an Engine` |

```java
// Inheritance
class Dog extends Animal { }

// Composition (preferred for flexibility)
class Car {
    Engine engine; // Car HAS an Engine

    Car() { this.engine = new Engine(); }
}
```

**Interview tip:** "Favor composition over inheritance" is a well-known design principle. Inheritance breaks encapsulation (child depends on parent internals); composition is more modular.

---

## 12. SOLID Principles

| Letter | Principle | Meaning |
|---|---|---|
| S | Single Responsibility | A class should have only one reason to change |
| O | Open/Closed | Open for extension, closed for modification |
| L | Liskov Substitution | Subtypes must be substitutable for their base types |
| I | Interface Segregation | Don't force classes to implement methods they don't use |
| D | Dependency Inversion | Depend on abstractions, not concrete implementations |

**Quick examples:**

**S** – Don't put DB logic + email logic in one class. Split them.

**O** – Use polymorphism: add new behavior by creating new subclasses, not editing existing code.

**L** – If `Bird` has a `fly()` method, and `Penguin extends Bird`, calling `penguin.fly()` breaks LSP. Fix by redesigning the hierarchy.

**I** – Don't put `print()`, `scan()`, `fax()` all in one `IMachine` interface. Split into `IPrinter`, `IScanner`.

**D** – A `Switch` class shouldn't depend on a `LightBulb` class directly. It should depend on an `IDevice` interface.

---

## 13. Common Design Patterns (OOP-based)

### Creational
- **Singleton** – Only one instance. Used for DB connections, config.
```java
class Singleton {
    private static Singleton instance;
    private Singleton() {}
    public static Singleton getInstance() {
        if (instance == null) instance = new Singleton();
        return instance;
    }
}
```
- **Factory** – Creates objects without exposing instantiation logic.
- **Builder** – Constructs complex objects step by step.

### Structural
- **Decorator** – Adds behavior to objects dynamically without subclassing.
- **Adapter** – Makes incompatible interfaces work together.

### Behavioral
- **Observer** – One-to-many dependency; when one object changes, all dependents are notified. (Event listeners)
- **Strategy** – Define a family of algorithms, encapsulate each, make them interchangeable.

---

## 14. Coupling and Cohesion

- **Coupling** – degree of dependency between classes. **Low coupling = good.** Classes should be as independent as possible.
- **Cohesion** – degree to which elements inside a class belong together. **High cohesion = good.** A class should do one thing well.

---

## 15. Association, Aggregation, Composition

| Concept | Description | Lifecycle Dependency |
|---|---|---|
| Association | General relationship between two classes | Independent |
| Aggregation | "has-a" — child can exist without parent | Independent |
| Composition | Strong "has-a" — child cannot exist without parent | Dependent |

```
Association:  Teacher ↔ Student (both exist independently)
Aggregation:  Department has Professors (Professor exists without Department)
Composition:  House has Rooms (Room cannot exist without House)
```

---

## 16. Shallow Copy vs Deep Copy

- **Shallow copy** – copies object references (both copies point to same nested objects)
- **Deep copy** – copies everything recursively (completely independent copy)

```java
// Shallow — changing address in copy affects original
Person copy = original.clone();

// Deep — must manually clone nested objects
Person deepCopy = new Person(original.name, new Address(original.address.city));
```

---

## 17. Quick-Fire Interview Questions & Answers

**Q: Can we override a static method?**
No. Static methods are resolved at compile time (class-level). They can be *hidden*, not overridden.

**Q: Can a constructor be private?**
Yes. Used in Singleton pattern to prevent external instantiation.

**Q: Can an abstract class have a constructor?**
Yes. It's called when a subclass object is created via `super()`.

**Q: Can we instantiate an abstract class?**
No. But we can have a reference of abstract class type pointing to a subclass object.

**Q: What is an interface marker?**
An empty interface with no methods, used to signal metadata (e.g., `Serializable`, `Cloneable` in Java).

**Q: Difference between == and equals()?**
`==` compares references (memory address). `equals()` compares values (logical equality).

**Q: What is method hiding?**
When a subclass defines a static method with the same signature as a parent static method. The parent method is *hidden*, not overridden.

**Q: What is the difference between early binding and late binding?**
Early binding (compile-time) = method overloading. Late binding (runtime) = method overriding via virtual dispatch.

**Q: Is Java purely object-oriented?**
No — because of primitive types (`int`, `char`, etc.) and static methods/variables which don't belong to objects.

---

## 18. Summary Cheat Sheet

| Concept | One-liner |
|---|---|
| Encapsulation | Hide data, expose behavior via methods |
| Abstraction | Hide how, show what |
| Inheritance | Reuse parent's code; "is-a" |
| Polymorphism | Same interface, different behavior |
| Overloading | Same name, different params, same class |
| Overriding | Same signature, child redefines parent |
| Abstract class | Partial implementation, can't instantiate |
| Interface | Pure contract, multiple implementable |
| Composition | "has-a", prefer over inheritance |
| Singleton | One instance, global access point |
| Coupling | Lower = better |
| Cohesion | Higher = better |