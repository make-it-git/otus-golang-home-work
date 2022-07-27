`go test -bench=. -benchtime=10s > old_10s.txt`

`go test -bench=. -benchtime=10s > new_10s.txt`

`benchstat old_10s.txt new_10s.txt > benchstat.txt`

Результаты оптимизации доступны в benchstat.txt
