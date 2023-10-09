# Test task for Zero Agency
description:
Create a JSON REST server with two endpoints:

POST /edit/:Id - Modify a news item by Id
GET /list - List of news items
Use MySQL for the database (PostgreSQL is also acceptable as we have migrated to it).

Use Fiber as the server framework. Utilize Reform for working with the database.

Establish a database connection using a connection pool. Configure all settings through environment variables and/or Viper.

Input data format for the first endpoint:
```
{
  "Id": 64,
  "Title": "Lorem ipsum",
  "Content": "Dolor sit amet <b>foo</b>",
  "Categories": [1,2,3]
}
```
If any of the fields are not provided, that field should not be updated.

Output data format for the list endpoint:
```
{
  "Success": true,
  "News": [
    {
      "Id": 64,
      "Title": "Lorem ipsum",
      "Content": "Dolor sit amet <b>foo</b>",
      "Categories": [1,2,3]
    },
    {
      "Id": 1,
      "Title": "first",
      "Content": "tratata",
      "Categories": [1]
      }
    ]
}
```
Requirements and Preferences:

- If you are familiar with Docker, it would be preferable to see the service packaged in a container.
- Additionally, it would be appreciated to have (bonus points compared to other candidates):
    - Authorization through the Authorization header and proper code structuring and grouping of endpoints.
    - Field validation during editing.
    - Pagination on the list endpoint.
    - Proper logging using any popular logger (e.g., logrus).
    - Adequate error handling.


# Comments for checking:

Docker:
  - docker-compose up --build
  (*server may not start from the first time due to postgreSQL needs some time to create 'zero_agency_db')

Postman collection is located in root repository:
- [https://github.com/ERupshis/effective_mobile/blob/main/Effective%20Mobile%20test.postman_collection.json](https://github.com/ERupshis/zero_agency_test/blob/master/zero%20agency.postman_collection.json)

