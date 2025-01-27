# Distributed Data Aggregation system - ControllerApp

The controller acts as an application server, mediating between the client and the execution module responsible for performing computations. The core functionality of the controller lies in managing communication across two layers: with the client and with the execution module. The first responsibility is achieved through a public web API that handles user queries. The second layer comes into play during query processing, where the controller distributes tasks among all available executors and collects results from one of them.

The controller's public API enables the following functionalities:

- **Executing Aggregate Queries:**  
  Users can perform data aggregation with functions such as `average`, `sum`, `minimum`, `maximum`, and `count` on grouped data. From the uploaded table, users select grouping columns and aggregated columns, and specify the aggregation function for each aggregated column.

- **Retrieving the List of Available Tables:**  
  Provides a list of table names, corresponding to folder names in the main data directory.

- **Retrieving the List of Columns for a Specific Table:**  
  For a given table, the API returns the list of columns available for grouping and aggregation. The table schema is read from the metadata of the `.parquet` file located in the folder named after the table. Columns are filtered to include only those supported by the system, with a distinction between columns usable for grouping and aggregation.

- **Uploading a `.parquet` File to a Table:**  
  Users can upload a `.parquet` file to a table with a specified name. If a folder with the specified name does not exist in the main data directory, it is created. For existing tables, the system validates the request to ensure the file schema matches the table schema. If the schemas differ, an error message is returned. Valid files are stored in the folder corresponding to the specified table name.


## Prerequisites
Make sure you have Go (to run the app) and Docker (to create Docker image) installed.

## Run locally 
In the project directory, create a `.env` file with the required configuration settings, eg.
```
EXECUTOR_ADDRESSES="localhost:8080"
CONTROLLER_PORT=":3000"
DATA_PATH="/home/user/Desktop/data"
MAIN_EXECUTOR_IDX=-0
SWAGGER_HOST="localhost:3000"
ALLOWED_ORIGIN="http://localhost:3006"
EXECUTOR_EXECUTOR_PORT=8081
```
Run the program using:
   ``` bash
   $ go run controller
   ```
Use swagger interface and send a request in the correct format to the already running EXECUTOR.

## Docker image 

To create Docker image run:
   ``` bash
   $ sudo sh build_image.sh
   ```
