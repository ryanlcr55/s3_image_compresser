## Dependency

1. [libvips](https://www.libvips.org/)

## Setting

### aws credential
1. create file 
    ```
    touche ~/.aws/credentials
    ```

2.  edit content

    ```
    [default]
    aws_access_key_id = <AWS_ACCESS_KEY_ID>
    aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>
    ```
### Project Environment
1. copy file 

    ```
    cp .env.example .env
    ```

2. env list

Name          | Description  
--------------|:-----:
REGION        | bucket region
BUCKET_NAME   | bucket name
PREFIX        | bucket prefix
BOUNDARY_SIZE | compress file size greater than this value(KB) 


## Execute

```
go run cmd/compress/main.go
```