sqlite3 /tmp/employees.db "drop table employee;" > /dev/null 2>&1
sqlite3 /tmp/employees.db "create table employee(id INTEGER PRIMARY KEY, name, email, phone, dob, salary NUMERIC); insert into employee values (1, 'Alice', 'alice@bigco.rp', '202-555-5555', '04-01-1956', '75000'), (2, 'Bob', 'bob@bigco.rp', '323-867-5309', '12-31-1984', '40000');" > /dev/null 2>&1