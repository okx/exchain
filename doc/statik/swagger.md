#MAC下的swagger使用

###安装go-swagger:
```
//brew install go-swagger的版本会报错
go get -u -v github.com/go-swagger/go-swagger
cd $GOPATH/src/go-swagger/go-swagger/cmd/swagger
go install .
```

###安装statik:
`go get github.com/rakyll/statik`

为了方便前后端接口定义与调试，使用swagger来自动生成接口定义文档。生成步骤如下：
1. 在与rest.go同目录下的rest_doc.go中添加对应的接口描述；
2. cd cmd/okchaincli
3. swagger generate spec -o dex.yaml
4. mv dex.yaml ../../doc/swagger-ui/dex/
5. 在okdex目录下执行restart.sh
6. 运行okdexcli rest-server --chain-id=okchain
7. 在浏览器中打开http://localhost:1317/swagger-ui/dex/indext.html
在显示的页面中可以对各接口指定参数，执行并返回结果。