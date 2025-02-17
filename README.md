# SGNL Adapter - PagerDuty

SGNL Adapter for PagerDuty SoR

## Background Information

### Terminology

**Adapter** - A simple gRPC server that queries an external API and parses the response into a format suitable for the SGNL ingestion service. More information on adapters can be found in the [What is an adapter?](#what-is-an-adapter?) section.

**SGNL ingestion service** - One of SGNL's core microservices which is responsible for ingesting external data into SGNL's graph database.

**System of Record (SoR)** - An external system that provides data to be ingested into SGNL, typically via an API. For example, a CRM tool like Salesforce or HubSpot.

**Entity** - A SGNL term to represent a type of data within an SoR. An entity typically maps to a single API resource. For example, the IdentityNow Account entity is retrieved from the `/accounts` endpoint.

### What is an adapter?

An adapter is a gRPC server which has two main responsibilities:

1. Making requests to a System of Record (SoR) to retrieve data.
2. Transforming the response from the SoR into a format that can be consumed by the SGNL ingestion service.

An adapter is **stateless**. It simply acts as a proxy to send requests to SoRs and parse the responses. An adapter uses SGNL's [adapter-framework](https://github.com/SGNL-ai/adapter-framework) under the hood.

Requests to the adapter invoke the `GetPage` method.

### Adapter Authentication

An adapter authenticates incoming gRPC requests via the `token` metadata key. The value of this key must match one of the tokens in the `ADAPTER_TOKENS` file that you define. The `ADAPTER_TOKENS` file is a JSON array of strings, where each string is a token. Only one token is required, but multiple tokens can be defined in the event that a token needs to be rotated.

For example, an `ADAPTER_TOKENS` file may look like:

```
["<token1>", "<token2>", ...]
```

Once this file is created, set the `AUTH_TOKENS_PATH` environment variable to the path of the `ADAPTER_TOKENS` file. More information on starting an adapter is discussed below in the [Getting Started](#1-getting-started) section.

## PagerDuty Adapter

## Running the adapter


1. Create a JSON file (for example, `authTokens.json`) that will contain the tokens used to authenticate requests to the adapter. The tokens must be stored in the following format and note down the path of the file. 

   ```
   ["<token1>", "<token2>", ...]
   ```

1. If you don't need to build a Docker image, you can directly run the adapter. Set the `AUTH_TOKENS_PATH` environment variable to the path of the tokens file created in the previous step. Then run `go run cmd/adapter/main.go`. Otherwise, proceed to the next step.

1. Build the Docker image with the `adapter` command.
   ```
   docker build -t adapter:latest .
   ```
   **WARNING:** The image will contain the `ADAPTER_TOKENS` secrets file. **Do not push** this image to a public registry.
1. Run the adapter as a Docker container.
   ```
   docker run --rm -it -e AUTH_TOKENS_PATH=/path/to/file adapter:latest
   ```

### 2. PagerDuty SoR

More information on the PagerDuty API can be found [here](https://developer.pagerduty.com/api-reference/e65c5833eeb07-pager-duty-api).

#### Entity Endpoints

Supports the following endpoints:

1. teams
2. users
3. vendors
```bash
https://api.pagerduty.com/teams # endpoint
```

Add more endpoints as needed in the `datasource.go` file and update the `ValidEntityExternalIDs` map.


#### Response Schemas

The response schemas might vary for each entity. For example, a teams entity response may look like

```jsonc
{
  "teams": [
    {
      "id": "PQZPQGI",
      "name": "North American Space Agency (NASA)",
      "description": null,
      "type": "team",
      "summary": "North American Space Agency (NASA)",
      "self": "https://api.pagerduty.com/teams/PQZPQGI",
      "html_url": "https://pdt-apidocs.pagerduty.com/teams/PQZPQGI",
      "default_role": "manager",
      "parent": null
    }
  ],
  "limit": 1,
  "offset": 0,
  "total": 2,
  "more": true
}
```

#### Authentication

Test Instance of PagerDuty API supports authentication via API token

Basic Auth credentials and Bearer tokens are passed directly to an adapter in a `GetPage` request.

#### API Restrictions

The request restrictions for each entity. For example,

--**limit**: For Pagerduty paginated APIs, this is the maximum number of results that can be returned in a single request. Corresponds to the `PageSize` field in the `Request` object.
- **offset.** For Pagerduty paginated APIs, this is the number of results to skip before returning the next set of results. Corresponds to the `Cursor` field in the `Request` object.


### 3. Understanding the Adapter 

A simplified flow chart of an incoming gRPC request to an adapter is shown below:

![Adapter Flow](docs/assets/adapter_flow.png)

1. To Add more entities, update the `ValidEntityExternalIDs` map in the `datasource.go` file.

### 4. Local Testing

As specified in the [Getting Started](#1-getting-started) section, you can run the adapter locally either through Docker or directly with `go run`.

```go
go run cmd/adapter/main.go
```

By default, the adapter will listen on port 8080.

Using Postman, you can send a gRPC request to the adapter.

1. Define the [`GetPage` Protobuf definition](https://github.com/SGNL-ai/adapter-framework/blob/f2cafb0d963b54c350350967906ce59776d720a1/api/adapter/v1/adapter.proto).

![Define the `GetPage` Protobuf definition](/docs/assets/postman_proto_definition.png)

2. In the sidebar, click on **Collections** and create a new collection with the type set to **gRPC**.

3. Within this new collection, create a new gRPC request. Enter the URL of the adapter (e.g. `http://localhost:8080`) and select the `GetPage` method, which should be available in the dropdown if step 1 was completed successfully.

![Create a new gRPC request](/docs/assets/postman_new_grpc_request.png)

4. In the **Metadata** tab, add a `token` key and set the value to one of the tokens in the `ADAPTER_TOKENS` file.

5. In the **Message** tab, enter the `GetPage` request. It must follow the schema defined in step 1.

An example gRPC request to Fetch all teams:

```jsonc
{
    "cursor": "",
    "datasource": {
        "type": "Test-1.0.0",
        "address": "https://api.pagerduty.com",
        "auth": {
            "http_authorization": "Token token=y_NbAkKc66ryYTWUXYEu"
        },
        "config": "e30=",
        "id": "Test"
    },
    "entity": {
        "attributes": [
            {
                "external_id": "id",
                "type": "ATTRIBUTE_TYPE_STRING",
                "id": "id"
            }
        ],
        "external_id": "teams",
        "id": "Team",
        "ordered": false
    },
    "page_size": "1",
    "total": true
}
```

Response

```jsonc
{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "values": [
                            {
                                "string_value": "PQZPQGI"
                            }
                        ],
                        "id": "id"
                    }
                ],
                "child_objects": []
            }
        ],
        "next_cursor": "1"
    }
}
```


