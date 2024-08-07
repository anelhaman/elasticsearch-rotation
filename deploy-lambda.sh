#!/bin/bash -x

aws lambda update-function-code \
    --function-name <your-lambda-name> \
    --zip-file fileb://${PWD}/deployment.zip >> /dev/null

# # option to set longer timeout
# aws lambda update-function-configuration \
#    --function-name <your-lambda-name> --timeout 10