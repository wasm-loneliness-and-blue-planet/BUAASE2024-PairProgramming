cd t3-1-go
tinygo build -o main.wasm -target wasm ./main.go
cd ../t3-2-go
tinygo build -o main.wasm -target wasm -opt 2 -gc leaking ./main.go
cd ..
