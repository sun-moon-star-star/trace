```
+-------------+------------------+------+-----+----------------------+-------------------+
| Field       | Type             | Null | Key | Default              | Extra             |
+-------------+------------------+------+-----+----------------------+-------------------+
| id          | int(11) unsigned | NO   | PRI | NULL                 | auto_increment    |
| start_time  | datetime(6)      | NO   |     | CURRENT_TIMESTAMP(6) | DEFAULT_GENERATED |
| update_time | datetime(6)      | NO   |     | CURRENT_TIMESTAMP(6) | DEFAULT_GENERATED |
| used_value  | int(4) unsigned  | NO   |     | NULL                 |                   |
| identifier  | varchar(255)     | YES  |     | NULL                 |                   |
+-------------+------------------+------+-----+----------------------+-------------------+
```

1. 注册project_id, 有效时间为30min, 如果使用需要定时更新
2. 每次分配的时候查看历史identifier是否有相匹配的值，如果有且更新时间超过30min，则使用
