Understood. Here is the plain text content that you can copy directly into a file named `api.md`.

````
# Checkout API Documentation

This document provides details on the available endpoints for the Checkout API.

## API Server Base URL
`http://localhost:8080`

---

## 1. Create a new Checkout Session

Initializes a new, empty checkout session and returns a unique ID for it. This ID must be used in subsequent requests to scan items and get the total price.

*   **Endpoint**: `POST /checkouts`
*   **Method**: `POST`

### Request Body
No request body is required.

### Responses

#### ✅ **Success: 201 Created**
Returned when a checkout session is successfully created. The body contains the unique ID for the new session.

**Response Body:**
```json
{
  "checkoutId": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
````

#### ❌ **Error: 500 Internal Server Error**

Returned if the server fails to save the new session.

**Response Body:**

```json
{
  "error": "could not save session"
}
```

---

## 2. Scan an Item

Scans an item and adds it to an existing checkout session. The server will recalculate the total price based on its pricing rules.

- **Endpoint**: `POST /checkouts/{checkoutID}/scan`
- **Method**: `POST`

### Path Parameters

| Parameter    | Type   | Description                                          |
| :----------- | :----- | :--------------------------------------------------- |
| `checkoutID` | string | **Required**. The unique ID of the checkout session. |

### Request Body

The request body must be a JSON object containing the SKU of the item to scan.

**Body:**

```json
{
  "sku": "A"
}
```

### Responses

#### ✅ **Success: 204 No Content**

Returned when the item is successfully scanned and added to the session. No response body is returned.

#### ❌ **Error: 400 Bad Request**

Returned if the request body is invalid or if the provided SKU does not exist in the pricing rules.

**Response Body (Example: Invalid SKU):**

```json
{
  "error": "invalid sku: vga"
}
```

**Response Body (Example: Malformed JSON):**

```json
{
  "error": "invalid request body"
}
```

#### ❌ **Error: 404 Not Found**

Returned if no session exists for the given `checkoutID`.

**Response Body:**

```json
{
  "error": "session not found"
}
```

#### ❌ **Error: 500 Internal Server Error**

Returned if the server fails to save the session after the scan.

**Response Body:**

```json
{
  "error": "could not save session"
}
```

---

## 3. Get Total Price

Retrieves the current total price for all items scanned in a specific checkout session.

- **Endpoint**: `GET /checkouts/{checkoutID}`
- **Method**: `GET`

### Path Parameters

| Parameter    | Type   | Description                                          |
| :----------- | :----- | :--------------------------------------------------- |
| `checkoutID` | string | **Required**. The unique ID of the checkout session. |

### Request Body

No request body is required.

### Responses

#### ✅ **Success: 200 OK**

Returned when the total price is successfully retrieved. The price is an integer representing the total in the smallest currency unit (e.g., cents).

**Response Body:**

```json
{
  "checkoutId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "totalPrice": 205
}
```

#### ❌ **Error: 404 Not Found**

Returned if no session exists for the given `checkoutID`.

**Response Body:**

```json
{
  "error": "session not found"
}
```

```

```
