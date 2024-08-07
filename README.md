
## Environment variables

1. set INDEX_AGE_LIMIT_DAYS to 90 to delete indices older than 90 days.
2. Live Mode: Set DRY_RUN=false (or remove it) to allow actual deletions.
3. set ES_URL: The URL of your Elasticsearch cluster (e.g., https://your-elasticsearch-domain:9200).
4. set ES_USERNAME: (Optional) The Elasticsearch username if authentication is required.
5. set ES_PASSWORD: (Optional) The Elasticsearch password if authentication is required.


---
# Elasticsearch Index Rotation Lambda
## Overview
This repository contains an AWS Lambda function designed to manage Elasticsearch index rotation. The function automates the deletion of old indices, retaining only a specified number of the most recent indices. This helps manage storage costs and maintain optimal performance in Elasticsearch.

## Features
Automated Index Rotation: Automatically deletes old Elasticsearch indices based on a retention policy.
Configurable Retention: Allows you to specify the number of indices to keep through environment variables.

Dry Run Mode: Includes a dry run mode to preview the indices that would be deleted without performing any actual deletions.

Environment-Specific Configurations: Supports different configurations for non-production and production environments.

## Prerequisites
AWS account with permissions to create and manage Lambda functions.

Elasticsearch cluster accessible from the Lambda function.

AWS CLI configured on your development machine.
Configuration

## Environment variables

1. set INDEX_AGE_LIMIT_DAYS to 90 to delete indices older than 90 days. (any number following your policy)
2. Live Mode: Set DRY_RUN=false (or remove it) to allow actual deletions.
3. set ES_URL: The URL of your Elasticsearch cluster (e.g., https://your-elasticsearch-domain:9200 or https://your-opensearch-domain:443).
4. set ES_USERNAME: (Optional) The Elasticsearch username if authentication is required.
5. set ES_PASSWORD: (Optional) The Elasticsearch password if authentication is required.

## Setup and Deployment
Clone the repository:
bash

```
git clone https://github.com/anelhaman/elasticsearch-rotation.git
cd elasticsearch-rotation
```
```
$ ./build.sh
```
### Optional
```
$ ./deploy-lambda.sh
```


## Usage
Once deployed, the Lambda function will execute based on the configured schedule (e.g., via CloudWatch Events) for example cron(0 0 ? * 1 *) as weekly 

The function will evaluate the indices in your Elasticsearch cluster, retaining only the specified number of the most recent indices, and deleting the rest.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License
This project is licensed under the MIT License.