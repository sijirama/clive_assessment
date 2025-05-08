
# Insight Invoice – Candidate TODO

This file is intentionally light.  
**Replace every section marked `TODO` with your own content** while completing the assessment.

## Project Purpose (TODO)
Explain what this project does in 1–2 sentences.

This project processes invoices to generate periodic financial insight reports using a cron job and Google Gemini LLM, exposing APIs to create, retrieve, and analyze invoices and their insights.

## Quick Start (TODO)
Provide commands to spin up the stack from scratch.
1. fill up all the .env.example by creating the .env equivalent in their respective directories
  - .env.example in ./worker
  - .env.example in ./
2. Spin up the stack using Docker Compose:
```bash
docker-compose up --build
```
3. the golang worker runs in http://localhost:8080 while the express server runs at http://localhost:3000


## API Reference (TODO)
Document each endpoint you build.

### Get 8000/insight/report
Description: Retrieves the most recent InsightReport with metrics like total spend, largest vendor, overdue count, and LLM-generated recommendations.

Response:

- 200 OK: Returns the latest report.

```json
{
  "report": {
    "ID": "f087ff76-aa94-40a3-8577-9c6d73ba2509",
    "TotalSpend": 2501.5,
    "LargestVendor": "Siji Corporation",
    "OverdueCount": 2,
    "Anomalies": {"anomalies": []},
    "CostSavingRecommendation": "Negotiate bulk discounts...",
    "ReportDate": "2025-05-08T14:19:40.036Z",
    "CreatedAt": "2025-05-08T14:19:43.722Z",
    "UpdatedAt": "2025-05-08T14:19:43.722Z"
  }
}
```

- 404 Not Found: No reports available.
- 500 Internal Server Error: Database or server error.

### POST 3000/
Description: Creates a new invoice.
Request Body:
```json
{
  "vendor": "Siji Corporation",
  "amount": 1250.75,
  "dueDate": "2025-06-15",
  "fileName": "invoice_acme_may2025.pdf",
  "rawInvoice": {
    "invoiceNumber": "INV-2025-0423",
    "lineItems": [
      {
        "description": "Web Development Services",
        "quantity": 40,
        "unitPrice": 25.00,
        "total": 1000.00
      },
      {
        "description": "Server Maintenance",
        "quantity": 5,
        "unitPrice": 50.15,
        "total": 250.75
      }
    ],
    "taxRate": 0,
    "notes": "Payment due within 30 days"
  },
  "isPaid": false
}
```

Response:

- 201 Created: Returns the created invoice.
```json
{
  "success": true,
  "message": "Invoice created successfully",
  "data": {
    "vendor": "Siji Corporation",
    "amount": 1250.75,
    "dueDate": "2025-06-15T00:00:00.000Z",
    "fileName": "invoice_acme_may2025.pdf",
    "rawInvoice": {
      "invoiceNumber": "INV-2025-0423",
      "lineItems": [
        {
          "description": "Web Development Services",
          "quantity": 40,
          "unitPrice": 25,
          "total": 1000
        },
        {
          "description": "Server Maintenance",
          "quantity": 5,
          "unitPrice": 50.15,
          "total": 250.75
        }
      ],
      "taxRate": 0,
      "notes": "Payment due within 30 days"
    },
    "isPaid": false,
    "paidDate": null,
    "id": "3c0d232a-e3d9-4df0-be2f-76b8906c465f",
    "createdAt": "2025-05-08T12:24:24.123Z",
    "updatedAt": "2025-05-08T12:24:24.123Z"
  }
}
```

- 400 Bad Request: Invalid input.
- 500 Internal Server Error: Database error.

### GET /invoices/:id
Description: Retrieves an invoice by its ID.
Parameters:
id: UUID of the invoice (e.g., 123e4567-e89b-12d3-a456-426614174000).

Response:
- 200 OK: Returns the invoice.
```json
{
  "success": true,
  "data": {
    "id": "3c0d232a-e3d9-4df0-be2f-76b8906c465f",
    "vendor": "Siji Corporation",
    "amount": "1250.75",
    "dueDate": "2025-06-15",
    "fileName": "invoice_acme_may2025.pdf",
    "rawInvoice": {
      "invoiceNumber": "INV-2025-0423",
      "lineItems": [
        {
          "description": "Web Development Services",
          "quantity": 40,
          "unitPrice": 25,
          "total": 1000
        },
        {
          "description": "Server Maintenance",
          "quantity": 5,
          "unitPrice": 50.15,
          "total": 250.75
        }
      ],
      "taxRate": 0,
      "notes": "Payment due within 30 days"
    },
    "isPaid": false,
    "paidDate": null,
    "createdAt": "2025-05-08T12:24:24.123Z",
    "updatedAt": "2025-05-08T12:24:24.123Z"
  },
  "document": [
    {
      "externalId": "3c0d232a-e3d9-4df0-be2f-76b8906c465f",
      "vendor": "Siji Corporation",
      "amount": 1250.75,
      "dueDate": "2025-06-15T00:00:00.000Z",
      "fileName": "invoice_acme_may2025.pdf",
      "rawInvoice": {
        "invoiceNumber": "INV-2025-0423",
        "lineItems": [
          {
            "description": "Web Development Services",
            "quantity": 40,
            "unitPrice": 25,
            "total": 1000
          },
          {
            "description": "Server Maintenance",
            "quantity": 5,
            "unitPrice": 50.15,
            "total": 250.75
          }
        ],
        "taxRate": 0,
        "notes": "Payment due within 30 days"
      },
      "isPaid": false,
      "paidDate": null,
      "createdAt": "2025-05-08T12:24:24.170Z",
      "updatedAt": "2025-05-08T12:24:24.170Z",
      "__v": 0,
      "id": "681ca278f943c678f9d3a8dd"
    }
  ]
}
```

- 404 Not Found: Invoice not found.
- 500 Internal Server Error: Database error.

## Design Notes (TODO)
Briefly explain architecture choices, trade‑offs, and how you used AI assistance.

### Architecture: 
- Built with Go, Gin for the HTTP server, GORM for PostgreSQL ORM, and robfig/cron for scheduling. A cron job runs SummariseAndSave to process invoices, detect anomalies, and generate LLM-based cost-saving recommendations using Google Gemini (gemini-1.5-pro).
- Express for the Node server with TypeORM for the Postgres ORM and Mongoose for the Mongo ODM

### Trade-offs:
- Used GORM for simplicity, but raw SQL could optimize complex queries.



---
