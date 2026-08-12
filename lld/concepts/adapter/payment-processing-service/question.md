
## Problem Statement

You are working on a Payment Processing Service for an e-commerce platform.

Your system already supports payments via Razorpay using the following interface:

```java
interface PaymentGateway {
    void pay(int amount);
}
```


The existing implementation:
```java
class RazorpayGateway implements PaymentGateway {
    public void pay(int amount) {
        System.out.println("Paid ₹" + amount + " using Razorpay");
    }
}
```

## New Requirement

Your company now wants to integrate a third-party payment provider – PayU.

However, PayU provides a different API which you cannot modify:

```java
class PayUGateway {
    public void makePayment(double amountInRupees) {
        System.out.println("Paid ₹" + amountInRupees + " using PayU");
    }
}
```

## Constraints

- ❌ You cannot change the PaymentGateway interface
- ❌ You cannot modify the PayUGateway class
- ✅ The client code should work with both Razorpay and PayU without any change
- ✅ Follow SOLID principles