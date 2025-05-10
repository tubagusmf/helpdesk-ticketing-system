# Helpdesk Ticketing System

Helpdesk Ticketing System is a backend application built with **Golang** to manage support tickets such as issue reporting, customer inquiries, or technical support requests.

## 🔧 Technologies Used

- **Golang** – Main programming language.
- **PostgreSQL** – Relational database for storing core data (tickets, users, comments, etc.).
- **Redis** – Used for caching ticket data to improve performance.
- **Elasticsearch** – For fast searching and filtering of ticket history data.
- **RabbitMQ** – Message broker used to send email notifications when tickets are created or updated.
- **Echo Framework** – For building RESTful HTTP APIs.
- **sql-migrate** – For handling database migrations.

## ✨ Key Features

- Ticket CRUD operations
- Ticket History tracking
- Comments and attachments on tickets
- Email notifications via RabbitMQ
- Ticket history search using Elasticsearch
- Redis caching for better performance

## ⚙️ Getting Started

### 1. Run the Application

Make sure you have the following services running:
- PostgreSQL
- Redis
- Elasticsearch (default: `http://localhost:9200`)
- RabbitMQ (default: `amqp://guest:guest@localhost:5672/`)

To run the application:

```bash
go run main.go httpsrv

```sql-migration
sql-migrate up // To apply database migrations:
sql-migrate new name_of_table // To create a new migration file
sql-migrate down // To undo the last migration


