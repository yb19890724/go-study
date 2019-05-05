#### 分析

- Instrumentation 仪表盘
    - 跟踪计数、延迟、健康状况和其他的周期性的或针对每个请求信息的仪表盘化，才能被认为是“生产环境”完备的
- Logging 日志

#### 启动服务 stringsvc2 mian.go

- 请求地址: localhost:8080/uppercase
- 请求方式: post
- 请求参数：body
```json
{
	"s":"test demo "
}
```

- 响应:response
```gotemplate
{
    "v": "TEST DEMO "
}
```

#### 查看分析
```gotemplate
# 帮助 GO U GC U 持续时间 U 秒 GC调用持续时间的摘要。
# 类型 go_gc_启动_秒 总结
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 5
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.11.6"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 636664
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 636664
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.443189e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 116
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.234368e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 636664
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.488064e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.835008e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 2010
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 0
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6715648e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 2126
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6912
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 27208
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 32768
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.055619e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 393216
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 393216
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.1891192e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 8
# HELP my_group_string_service_count_result The result of each count method.
# TYPE my_group_string_service_count_result summary
my_group_string_service_count_result{quantile="0.5"} 10
my_group_string_service_count_result{quantile="0.9"} 10
my_group_string_service_count_result{quantile="0.99"} 10
my_group_string_service_count_result_sum 10
my_group_string_service_count_result_count 1
# HELP my_group_string_service_request_count Number of requests received.
# TYPE my_group_string_service_request_count counter
my_group_string_service_request_count{error="false",method="count"} 1
my_group_string_service_request_count{error="false",method="uppercase"} 2
# HELP my_group_string_service_request_latency_microseconds Total duration of requests in microseconds.
# TYPE my_group_string_service_request_latency_microseconds summary
my_group_string_service_request_latency_microseconds{error="false",method="count",quantile="0.5"} 0.000146712
my_group_string_service_request_latency_microseconds{error="false",method="count",quantile="0.9"} 0.000146712
my_group_string_service_request_latency_microseconds{error="false",method="count",quantile="0.99"} 0.000146712
my_group_string_service_request_latency_microseconds_sum{error="false",method="count"} 0.000146712
my_group_string_service_request_latency_microseconds_count{error="false",method="count"} 1
my_group_string_service_request_latency_microseconds{error="false",method="uppercase",quantile="0.5"} 3.8751e-05
my_group_string_service_request_latency_microseconds{error="false",method="uppercase",quantile="0.9"} 0.000124038
my_group_string_service_request_latency_microseconds{error="false",method="uppercase",quantile="0.99"} 0.000124038
my_group_string_service_request_latency_microseconds_sum{error="false",method="uppercase"} 0.00016278900000000002
my_group_string_service_request_latency_microseconds_count{error="false",method="uppercase"} 2
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0

```
