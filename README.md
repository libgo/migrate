migration file name format:
  version_module_action.sql

./migrate -d "root:ddg1208@tcp(192.168.10.191:3306)/dolphin" -p "./test" -m "admin.org"
