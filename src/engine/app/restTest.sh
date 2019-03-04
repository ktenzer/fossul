#/bin/sh
curl -f -k -H 'Content-Type: application/json' -XPOST -d '{"profile": "myprofile","backupName": "mybackup","appPlugin": "myapp", "storagePlugin": "mystorage", "backupRetentions": [{"policy": "daily", "retentionDays": 5},{"policy": "weekly", "retentionDays": 10}],"preAppQuiesceCmd": "/path/to/pre_quiesce","appQuiesceCmd": "/path/to/app_quiesce","postAppQuiesceCmd": "/path/to/post_quiesce","backupCreateCmd": "/path/to/backup_create","preAppUnquiesceCmd": "/path/to/pre_unquiesce","appUnquiesceCmd": "/path/to/unquiesce","postAppUnquiesceCmd": "path/to/post_unquiesce"}' http://localhost:8001/quiesce
