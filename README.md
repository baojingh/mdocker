# mdocker

# ref

```bash
https://zhuanlan.zhihu.com/p/360242317

https://www.cnblogs.com/sevck/p/15033922.html

https://www.reddit.com/r/docker/comments/9s8yaf/using_the_sdk_how_can_i_wait_for_a_container_exec/


https://stackoverflow.com/questions/54747796/golang-docker-api-cant-run-exec-on-a-container

https://monkeywie.cn/2019/07/19/docker-web-terminal/
https://github.com/bigheadbros/docker-web-terminal
```


```azure
python的工程
https://github.com/AShadowMan/docker-web-terminal/blob/master/utility/myDocker.py
```


```azure
https://github.com/jesseduffield/lazydocker
```


```
https://docs.docker.com/engine/api/v1.41/
```

```azure
{
    "read": "2023-05-10T06:31:45.845159796Z",
    "preread": "2023-05-10T06:31:44.842822079Z",
    "pids_stats": {
        "current": 18
    },
    "blkio_stats": {
        "io_service_bytes_recursive": [
            {
                "major": 253,
                "minor": 0,
                "op": "Read",
                "value": 614400
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Write",
                "value": 0
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Sync",
                "value": 0
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Async",
                "value": 614400
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Total",
                "value": 614400
            }
        ],
        "io_serviced_recursive": [
            {
                "major": 253,
                "minor": 0,
                "op": "Read",
                "value": 17
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Write",
                "value": 0
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Sync",
                "value": 0
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Async",
                "value": 17
            },
            {
                "major": 253,
                "minor": 0,
                "op": "Total",
                "value": 17
            }
        ],
        "io_queue_recursive": [],
        "io_service_time_recursive": [],
        "io_wait_time_recursive": [],
        "io_merged_recursive": [],
        "io_time_recursive": [],
        "sectors_recursive": []
    },
    "num_procs": 0,
    "storage_stats": {},
    "cpu_stats": {
        "cpu_usage": {
            "total_usage": 20686605138,
            "percpu_usage": [
                10355933926,
                10330671212
            ],
            "usage_in_kernelmode": 10390000000,
            "usage_in_usermode": 10030000000
        },
        "system_cpu_usage": 2438625970000000,
        "online_cpus": 2,
        "throttling_data": {
            "periods": 0,
            "throttled_periods": 0,
            "throttled_time": 0
        }
    },
    "precpu_stats": {
        "cpu_usage": {
            "total_usage": 206
85447955,
            "percpu_usage": [
                10355933926,
                10329514029
            ],
            "usage_in_kernelmode": 10390000000,
            "usage_in_usermode": 10030000000
        },
        "system_cpu_usage": 2438623980000000,
        "online_cpus": 2,
        "throttling_data": {
            "periods": 0,
            "throttled_periods": 0,
            "throttled_time": 0
        }
    },
    "memory_stats": {
        "usage": 16027648,
        "max_usage": 17461248,
        "stats": {
            "active_anon": 15392768,
            "active_file": 356352,
            "cache": 634880,
            "dirty": 0,
            "hierarchical_memory_limit": 9223372036854771712,
            "hierarchical_memsw_limit": 9223372036854771712,
            "inactive_anon": 0,
            "inactive_file": 278528,
            "mapped_file": 294912,
            "pgfault": 47266,
            "pgmajfault": 4,
            "pgpgin": 18218,
            "pgpgout": 16349,
            "rss": 15392768,
            "rss_huge": 8388608,
            "total_active_anon": 15392768,
            "total_active_file": 356352,
            "total_cache": 634880,
            "total_dirty": 0,
            "total_inactive_anon": 0,
            "total_inactive_file": 278528,
            "total_mapped_file": 294912,
            "total_pgfault": 0,
            "total_pgmajfault": 0,
            "total_pgpgin": 0,
            "total_pgpgout": 0,
            "total_rss": 15392768,
            "total_rss_huge": 8388608,
            "total_unevictable": 0,
            "total_writeback": 0,
            "unevictable": 0,
            "writeback": 0
        },
        "limit": 3873652736
    },
    "name": "/redis",
    "id": "25edf5f9b9c4d1b56ef1f58d4e1c6e1f54151962a85e2af042ba49479efe528d"
}
```

