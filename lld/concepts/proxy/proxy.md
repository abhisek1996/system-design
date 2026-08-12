
**Intent**
- structural design pattern

- caching
- access
- pre processing/post processing
    - logging
    - event

## BFF
The Core Problem: Lender Migration
Your company is migrating from an old backend (loans-api) to a new backend (maestro). But you can't migrate all lenders at once—it happens lender by lender.

- Without proxy: The UI would need to know which backend to call for each lender.
- With proxy: UI calls one endpoint → BFF figures out the right backend.
## Proxy vs. Decorator (The "Intent")

While the structure (wrapping an interface) is identical, the intent is different:

| Pattern | Primary Intent | Relationship |
| :--- | :--- | :--- |
| **Proxy** | **Control**: Access control, caching, lazy initialization. Manages the lifecycle. | Usually **1-to-1** with the subject. |
| **Decorator** | **Enhance**: Adds new responsibilities or behaviors dynamically. | Often **dynamic** stacking. |

- **Proxy** *manages* the access to an object.
- **Decorator** *augments* the behavior of an object.
