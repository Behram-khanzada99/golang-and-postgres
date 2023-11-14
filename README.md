# golang-and-postgres
+ Inserting timestamps. <br />
+ Achieve concurrency in Golang through Go routines (pggoconcurrency) <br />
+ Inserted 1,000,000 (million) records while achieving concurrency <br />
Results are as followed. <br /><br />
**Benchmarks:** <br />
+ Inserted 100,000 records in 47,873,830.827 microseconds (47.8 seconds) in the Database and Read 100,000 records in 51,172.019 microseconds (0.051 seconds) (Go Language) ) <br />
+ Inserted 1,000,000 (million) records in 8 minutes and 41.788600694 seconds and Read 1,000,000 (million) records in 503,822.315 microseconds (0.5 seconds) (Go Language) ) <br /><br />
**Stress-test:** <br />
+ Performing concurrent CRUD operations seamlessly without encountering any errors. <br />
