# Wakyma Plugin – Reporting Service

This project is a lightweight **Go backend**(in the future there will be an interface) designed to integrate with the **Wakyma API** and provide a flexible framework for generating custom invoice reports.

It allows you to:
- Fetch invoices from Wakyma.
- Expose local HTTP endpoints for data retrieval.
- Build and extend multiple **report generators** (Excel).
- Add new report types without modifying core logic.

---

## Project Structure

```
backend/
├── entities/      # Shared domain models
├── usecases/
│   ├── listallinvoices/  # Core invoice data fetcher
│   └── createreport1/    # Daily report registration
├── main.go        # Server and route registration
└── Makefile
```

---

## Architecture

- **Entities** define shared data models (Invoice, Client, Pagination, etc.).
- **Usecases** contain business logic (fetching invoices or generating reports).
- **HTTP handlers** map endpoints to usecases.

Each report lives in its own folder under:

```
backend/usecases/<report-name>/
```

---

## Environment Variables

This file should be in the folder `backend` or in the same folder as the .exe file
```env
API_URL_WAKYMA=https://vets.wakyma.com/api/v3
API_KEY_WAKYMA=your_api_token
````

---

## Running

```bash
cd backend
make run
```

Server runs at:

```
http://localhost:8080
```

---

## Wakyma API

You can consult the documentation of the API in this link: [(https://vets.wakyma.com/public/api-docs/](https://vets.wakyma.com/public/api-docs/)

---

## Dev Commands

```bash
make fmt        # format code
make tidy       # tidy go.mod
make test       # run tests(no tests yet)
make lint       # lint (optional)
make clean      # remove build artifacts
```

---

## Notes

* Output files (`*.xlsx`, files inside /bin) are git-ignored.
* Keep API credentials private.
* .env file for the Go application must be in the backend folder

---

## License

Feel free to fork the repository and extend or modify its functionalities