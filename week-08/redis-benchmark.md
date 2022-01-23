作业内容：

1.使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

2.写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

# 0.实验前准备

[docker-compose 文件](docker-compose.yaml):利用 docker 在本地启动一个 Redis

## 作业1

### 1.1 执行命令
```shell
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 10 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 20 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 50 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 100 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 200 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 1000 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 5000 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 10000 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 15000 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 50000 --csv -t get,set >> txt.csv && \
redis-benchmark -h 127.0.0.1 -p 6378 -a pwd-turato -q -d 100000 --csv -t get,set >> txt.csv 
```

### 1.2 测试结果

| data_size_byte(单位,字节) | test | rps | avg_latency_ms | min_latency_ms |  | p50_latency_ms | p95_latency_ms | p99_latency_ms | max_latency_ms |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 10 | SET | 6645.84 | 7.413 | 0.92 |  | 6.479 | 14.695 | 22.751 | 103.359 |
| 10 | GET | 5420.94 | 9.079 | 0.952 |  | 7.575 | 19.487 | 33.535 | 146.431 |
| 20 | SET | 5613.88 | 8.763 | 0.992 |  | 7.423 | 18.319 | 29.327 | 151.551 |
| 20 | GET | 5735.26 | 8.575 | 0.808 |  | 7.247 | 17.823 | 30.991 | 127.359 |
| 50 | SET | 6583.28 | 7.476 | 0.904 |  | 6.551 | 14.855 | 23.215 | 125.887 |
| 50 | GET | 6995.45 | 7.034 | 0.944 |  | 6.183 | 13.575 | 20.319 | 142.975 |
| 100 | SET | 5664.12 | 8.712 | 0.88 |  | 7.327 | 18.415 | 30.927 | 146.943 |
| 100 | GET | 5626.83 | 8.789 | 1 |  | 7.415 | 18.559 | 30.191 | 175.615 |
| 200 | SET | 6459.95 | 7.628 | 0.816 |  | 6.607 | 15.247 | 24.175 | 137.727 |
| 200 | GET | 6837.61 | 7.188 | 0.816 |  | 6.335 | 13.855 | 21.343 | 110.911 |
| 1000 | SET | 6799.95 | 7.2 | 1.056 |  | 6.503 | 13.407 | 18.943 | 86.271 |
| 1000 | GET | 6824.54 | 7.213 | 0.864 |  | 6.319 | 14.183 | 21.919 | 122.815 |
| 5000 | SET | 5379.53 | 9.142 | 1.104 |  | 8.695 | 14.983 | 19.007 | 101.759 |
| 5000 | GET | 5260.11 | 9.441 | 0.92 |  | 7.895 | 18.639 | 34.495 | 165.375 |
| 10000 | SET | 3287.09 | 15.122 | 1.632 |  | 14.903 | 20.815 | 25.391 | 130.815 |
| 10000 | GET | 4509.58 | 11.04 | 1.344 |  | 9.111 | 17.583 | 44.543 | 1018.367 |
| 15000 | SET | 2406.62 | 20.711 | 1.672 |  | 20.207 | 27.807 | 34.815 | 140.031 |
| 15000 | GET | 3398.01 | 14.671 | 1.68 |  | 12.911 | 20.687 | 68.287 | 179.967 |
| 50000 | SET | 830.25 | 60.155 | 2.088 |  | 57.183 | 80.319 | 101.375 | 1100.799 |
| 50000 | GET | 1155.78 | 18.971 | 2.176 |  | 15.743 | 32.527 | 95.295 | 1029.119 |
| 100000 | SET | 435.18 | 114.781 | 13.528 |  | 108.799 | 140.543 | 184.319 | 1135.615 |
| 100000 | GET | 609.74 | 21.293 | 2.184 |  | 17.103 | 35.647 | 118.911 | 1029.119 |

### 1.3 结论

个人的结论：在本地测试的，没有考虑网络传输问题。
- value 数据越大影响吞吐：rps 指标表示每秒请求数，data_size_byte 越大，rps 逐渐下降，redis 吞吐逐渐下降。
- value 数据大小应控制在 5k 以内：数据大小在 5k 的之前，表现都差不多，超过5k之后，rps 和 latency 表现下滑

## 作业2

### 2.1 执行代码

该[测试代码](main.go)构造了一堆数据,这些数据的 value 尽可能不相同。
模拟写入 1W 条 key-value 数据, 其 value 大小在[10 20 50 100 200 1k 5k] 字节中。

### 2.2 测试结果

| data_size_byte | used_memory | used_memory_human | used_memory_rss | used_memory_rss_human | used_memory_peak | used_memory_peak_human | used_memory_peak_perc | used_memory_overhead | used_memory_startup | used_memory_dataset | used_memory_dataset_perc | allocator_allocated | allocator_active | allocator_resident | total_system_memory | total_system_memory_human | used_memory_lua | used_memory_lua_human | used_memory_scripts | used_memory_scripts_human | number_of_cached_scripts | maxmemory | maxmemory_human | maxmemory_policy | allocator_frag_ratio | allocator_frag_bytes | allocator_rss_ratio | allocator_rss_bytes | rss_overhead_ratio | rss_overhead_bytes | mem_fragmentation_ratio | mem_fragmentation_bytes | mem_not_counted_for_evict | mem_replication_backlog | mem_clients_slaves | mem_clients_normal | mem_aof_buffer | mem_allocator | active_defrag_running | lazyfree_pending_objects |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 10 | 1645728 | 1.57M | 57114624 | 54.47M | 52685728 | 50.25M | 3.12% | 1388984 | 791296 | 256744 | 30.05% | 1682512 | 1884160 | 55947264 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.12 | 201648 | 29.69 | 54063104 | 1.02 | 1167360 | 34.74 | 55470336 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 20 | 1725728 | 1.65M | 6512640 | 6.21M | 52685728 | 50.25M | 3.28% | 1388984 | 791296 | 336744 | 36.04% | 1758608 | 1957888 | 5349376 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.11 | 199280 | 2.73 | 3391488 | 1.22 | 1163264 | 3.78 | 4790280 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 50 | 2045728 | 1.95M | 6496256 | 6.20M | 52685728 | 50.25M | 3.88% | 1388984 | 791296 | 656744 | 52.35% | 2076592 | 2269184 | 5394432 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.09 | 192592 | 2.38 | 3125248 | 1.2 | 1101824 | 3.19 | 4456800 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 100 | 2605728 | 2.49M | 6979584 | 6.66M | 52685728 | 50.25M | 4.95% | 1388984 | 791296 | 1216744 | 67.06% | 2635360 | 2842624 | 5939200 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.08 | 207264 | 2.09 | 3096576 | 1.18 | 1040384 | 2.69 | 4386048 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 200 | 3725728 | 3.55M | 8273920 | 7.89M | 52685728 | 50.25M | 7.07% | 1388984 | 791296 | 2336744 | 79.63% | 3763024 | 3989504 | 7143424 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.06 | 226480 | 1.79 | 3153920 | 1.16 | 1130496 | 2.22 | 4548520 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 1000 | 11725728 | 11.18M | 16052224 | 15.31M | 52685728 | 50.25M | 22.26% | 1388984 | 791296 | 10336744 | 94.53% | 11741736 | 11948032 | 15020032 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1.02 | 206296 | 1.26 | 3072000 | 1.07 | 1032192 | 1.37 | 4351184 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |
| 5000 | 52685728 | 50.25M | 59113472 | 56.38M | 52685728 | 50.25M | 100.00% | 1388984 | 791296 | 51296744 | 98.85% | 52389208 | 52596736 | 58159104 | 2082197504 | 1.94G | 37888 | 37.00K | 0 | 0B | 0 | 0 | 0B | noeviction | 1 | 207528 | 1.11 | 5562368 | 1.02 | 954368 | 1.13 | 6799872 | 0 | 0 | 0 | 66616 | 0 | jemalloc-5.1.0 | 0 | 0 |

### 2.3 结论

单看 used_memory 这一指标来说： value 字节数 data_size 越大，平均每个 key-value，占用内存空间越大，但是其与 data_size_byte 并不是类似正比的关系。

比如 data_size_byte = 100, used_memory = 2.49M, 平均每个 key-value 占 2544KB,
data_size_byte = 5000, 数据大小扩大至 50 倍，使用空间仅增长至 20 倍。

所以结论：value 越大，存储结构越紧密。