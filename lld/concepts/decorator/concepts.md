
**Also known as: Wrapper**

**Intent**

- Attach additional responsibilities to an object dynamically. Decorators provide a flexible alternative to subclassing for extending functionality.
- To avoid class explosion
    - Base car, Base car with sunroof, Base car with leather seats, Base car with sunroof and leather seats : too many classes
    - Pizza with toppings : too many classes
    - Coffee machine with different options : too many classes

**Applicability**

Use the Decorator pattern when

- You need to add responsibilities to individual objects dynamically and transparently, without affecting the behavior of other objects from the same class.
- You want to avoid permanent extension of a class with new responsibilities.
- You need to add responsibilities to objects without modifying their structure.
- You want to support a combination of responsibilities that can be changed at runtime.

**Structure**

- Component: defines an interface for objects that can have responsibilities added to them dynamically.
- ConcreteComponent: a class representing the core functionality that can be decorated.
- Decorator: a class that implements the Component interface and holds a reference to a Component object.
- ConcreteDecorator: classes that add one or more responsibilities to the Component object.