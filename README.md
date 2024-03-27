Прокси для эндпоинта /metrics который на лету подсчитывает общее колличество метрик и их кардинальность.
Можно использовать как сайдкар перед любым приложением в кубере.

Позволит не нагружать Prometheus/VictoriaMetrics запросами вида, так же в случае высоконагруженных кластеров позволит
строить историчесские графики что в случае с запросом представленным ниже не представляется возможным для более менее
больших временных рядов:
```js
topk(10, sort_desc(count by (__name__) ({container="app", namespace="$namespace"})))
```

Для запуска требуется указать переменные окружения:
```bash
LISTEN_ADDR ":8080" - порт для запуска http сервера
METRICS_URL "http://localhost:58915/metrics" - источник метрик
```

При запросе на `http://$(LISTEN_ADDR)/metrics` отдаёт метрики с источника дополняя их расчётными значениями

Пример:
```bash
# HELP metrics_cardinality Metrics label cardinality
# TYPE metrics_cardinality gauge
metrics_cardinality{metric_name="go_gc_duration_seconds"} 1
metrics_cardinality{metric_name="redis_connections_open"} 1
metrics_cardinality{metric_name="go_memstats_lookups_total"} 1
metrics_cardinality{metric_name="go_info"} 1
metrics_cardinality{metric_name="grpc_server_handled_total"} 18
.....
    
# HELP metrics_count Count of metrics
# TYPE metrics_count gauge
metrics_count 122
```