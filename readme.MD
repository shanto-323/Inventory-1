# Inventory-01 
### A self-hosted inventory management service powered by Elasticsearch  
**Create, update, search, and analyze product data using fast full-text and structured queries.  
Designed to run locally or on your own server with full control over performance, data, and infrastructure.**


## Run Locally

Clone the project

```bash
  git clone git@github.com:shanto-323/Inventory-1.git
```

Go to the project directory

```bash
  cd project/path/Inventory-1
```

Start the server

```bash
  docker-compose up --build -d
```
***or***
```bash
  docker-compose down -v && docker-compose up --build 
```

## ***API Reference***
#### Object Structure

| Field     | Type               | Description                           |
|-----------|--------------------|-------------------------------------|
| `name`    | `string`           | Display name of the product          |
| `type`    | `string`           | Product category                     |
| `brand`   | `string`           | Manufacturer brand                   |
| `model`   | `string`           | Model number                        |
| `stock`   | `int` (integer)    | Inventory count                     |
| `specs`   | `map<string,string>`| Technical specifications (key-value pairs) |
| `warranty`| `string`           | Warranty duration                   |
| `supplier`| `string`           | Vendor or supplier name             |
| `note`    | `string`           | Internal notes or comments          |

```bash
{
    "name": "Ryzen 9",
    "type": "Processor",
    "brand": "AMD",
    "model": "5900X",
    "stock": 3,
    "specs": {
        "Cores": "12",
        "Threads": "24",
        "Base Clock": "3.7GHz"
    },
    "warranty": "3 years",
    "supplier": "AMD Distributors",
    "note": "High-end processor"
}
```

### Basic CRUD Operations
</br>

| Operation       | Method | Endpoint                                                                                   | Description                     |
|-----------------|--------|--------------------------------------------------------------------------------------------|---------------------------------|
| Create Product  | POST   | `https://localhost:8080/api/v1/products`                                                  | Add a new product               |
| Get Product     | GET    | `https://localhost:8080/api/v1/products/{id}`              | Get a single product by ID      |
| Update Product  | PUT    | `http://localhost:8080/api/v1/products/{id}`              | Update product by ID            |
| Delete Product  | DELETE | `https://localhost:8080/api/v1/products/{id}`              | Delete product by ID            |
</br></br>
### Analytics Functions

### 1. Stock Level Filter
    GET /api/v1/analytics/stock?level=any
This endpoint returns all products with stock <= level. The level query parameter is optional.
If not provided, the default baseline value is 3.</br>
**Example:**

    GET /api/v1/analytics/stock?level=10 → filters stock ≤ 10

    GET /api/v1/analytics/stock → filters stock ≤ 3 (default)

</br>

### 2. Advanced Product Search

    POST /api/v1/analytics/search
Performs an advanced product search using multiple optional fields in the request body.
If no fields are provided, it returns all products (acts like a "list all" endpoint).
| Field           | Type     | Description                                                  |
|------------------|----------|--------------------------------------------------------------|
| `search_string`  | `string` | Keyword to search across name, model, and other text fields  |
| `product_type`   | `string` | Filter by exact product type                                 |
| `product_brand`  | `string` | Filter by exact brand name                                   |
| `product_model`  | `string` | Filter by exact model number                                 |
| `min_stock`      | `int`    | Include products with stock greater than or equal to this    |
| `max_stock`      | `int`    | Include products with stock less than or equal to this       |
| `supplier`       | `string` | Filter by exact supplier/vendor name                         |

All fields in the request body are **optional**.  
If no fields are provided, the endpoint behaves like a **"Get All Products"** operation.

