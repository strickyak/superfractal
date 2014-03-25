set -ex
mkdir -p Keep

go run super_png/main.go -p='s;s' -n=10  -base Keep/1-SierByTriangs -mustwhite=1 -starttriang=1 -np=1

go run super_png/main.go -p='s;s' -n=10  -base Keep/2-SierBySquares -mustwhite=1  -np=1

go run super_png/main.go -p='W;X;Y;Z' -n=10  -base Keep/4-Monster  -np=24


go build main/coded/coded.go && for s in 1 2 3 4 5 6 7 8 9 ; do time ./coded  -o Keep/7-Fernpinski-$s.000.png -p 's;f' -n 1000000 -s $s -d 8 -fuzz=2 ; done

go build main/coded/coded.go && for s in 1 2 3 4 5 6 7 8 9 ; do time ./coded  -o Keep/8-SplicedFerns-$s.000.png -p 'f;A;B;C' -n 1000000 -s $s -d 8 -fuzz=2 ; done
