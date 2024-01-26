# Brief

- This is a Hypergro.ai Assignment.
- It is an example stock data collector and explorer (practically speaking)

## Live

- Web App: https://hypergroai.netlify.app
- APIs: https://hypergroaiassignment-production.up.railway.app

## Features

- API to Add/upload stock data in CSV file format to the server and migrate it into the database
- API to trigger migration of CSV file in stock data on the server
- API to get all the migrations logs.
- Ability to log in/signup with your Google account.
- API to Query for stocks based on certain options
  - Pagination is available, at `/stock?page=<PAGE>&size=<SIZE>`
  - Can/Will have to query with date using query params "date", example `/stock?date=<DATE>`. _note_: date is here in the format of "DD/MM/YYYY", if there is no data available for your given date, the latest data will be returned as default
  - Can search stock with full name or just with a few starting characters, example `/stock?s=<TEXT>`
  - Also get top stocks, that is stocks with sort in terms of gain descendingly, at `/stock/top`. _note_: this can/will support all the other query parameters.
- API to add stocks as user favorites, interact with them (read all), or update the user favorites list
- A ReactJS UI to showcase all of these features and APIs in a great UX. _note_: the API to upload a stock data CSV file and trigger a migration API can't be triggered with the React UI.

### Documentation

- Check the Postman collection [here](https://documenter.getpostman.com/view/26854281/2s9YyqhgpN)

## Setup

- Backend Setup

  1. CD into app (server)

  ```
      cd app
  ```

  2. Create ENV file

  ```
      cp sample.env .env
  ```

  3. Setup mongodbDB and Redis instance (Docker needed)

  ```
    docker-compose -f ../config/compose.yml up -d
  ```

  _note_: You can use other ways too if required.

  4. Setup Google Oauth Client and get credentials, refer to [here](https://support.google.com/cloud/answer/6158849?hl=en)
  5. Setup Google Cloud Service Account and give Cloud Storage Admin access, refer to [here](https://cloud.google.com/iam/docs/keys-create-delete)
  6. Add the json key of the service account into the env after converting it into string, and replace `\n` with `%n%`.
  7. Run server

  ```
      go run main
  ```

- Frontend Setup (optional)

  1. CD into client

  ```
  cd client
  ```

  2. Install packages

  ```
  yarn
  ```

  or `npm i`

  3. Create the "ENV" file.

  ```
  cp sample.env .env
  ```

  4. Run the server

  ```
  yarn dev
  ```

  or `npm run dev`
