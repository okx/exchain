# How to compile solidity and use abigen to generate go file

## 1. compile contracts

### bin file generation
```sh
 solc --bin ${your contract file} -o build
```

### abi file generation
```sh
 solc --abi ${your contract file} -o build
```

After previous two steps, there will be abi,bin file in **${cmd execute dir}/build** 

## 2. use geth tool abigen to generate go from abi file
```sh
abigen --bin=${your bin file} --abi=${your abi file} --pkg=main
```
> **tips**: if you abigen go file without bin file path, there will be no deploy method in go file.