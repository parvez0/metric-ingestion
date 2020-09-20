# metric-ingestion
Go backend application to record metric ingestion and give an aggregated value, it exposes 2 apis 
**/metrics** and **/report/** to push the data and fetch the result.

#### Usage
    docker pull grokkertech/metric-ingestion:v0.1.0
    docker run -d -e LOGLEVEL=info -p 8080:8080 -v $PWD/data:/data grokkertech/metric-ingestion:v0.1.0

This project uses a small SQLITE DB to store the data in real practise we should replace it with mysql DB.
The above docker commands starts a server on port 8080 and while initializing, it also creates a schema, for 
the table. 
 
### Push Metrics 
##### Request

    POST /metrics
```json
{
    "percentage_cpu_used": 30,
    "percentage_memory_used": 70
}    
```    
##### Response

    STATUS: 200
    BODY: metric recorded
     
##### Error
```json
{
    "Success": false,
    "Message": "Bad request",
    "data": {
        "message": "percentage_cpu_used and percentage_memory_used are required"
    }
}
``` 
OR    

    STATUS: 500
    BODY internal server error
    
### Fetch Report 

##### Request

Returns the records with max cpu or max memory consumption group by IP address of host machine

    GET /report?filter=cpu
    
```text
filter <string>{optional}: default to cpu, [cpu, memory]
```

##### Response
```json
[
    {
        "row_id": 0,
        "percentage_cpu_used": 60,
        "percentage_memory_used": 70,
        "ip": "::1",
        "date": "2020-09-20T17:44:15Z"
    }
]
```
##### Error
```json
{
    "Success": false,
    "Message": "Bad request",
    "data": {
        "message": "please provide a valid filter either cpu or memory"
    }
}
```