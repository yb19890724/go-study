### TLS认证


资料参考：https://www.jianshu.com/p/f109b9569a87

#### 证书制作

* 客户端和服务器通过gRPC相互通信，由于我们使用了protobuf来序列化和反序列化消息，因此消息的数据是二进制形式的。但是我们的通信是明文传输的，这在一些安全需求较高的场景中是不允许的，因此需要使用安全传输。幸好，在gRPC中，我们可以直接使用SSL/TLS，用来验证服务器，并对通信过程进行加密。

#### 生成私钥

- 进入keys这个目录，执行下面的命令来生成私钥：

```gotemplate
openssl genrsa -out server.key 2048
```
我们使用openssl genrsa命令来生成私钥，并用-out选项指定输出。最后一个参数2048表示的是生成密钥的位数，如果没有指定，那么默认就是512位。


#### 根据私钥生成CSR

如果想从一个认证中心（Certificate Authority，CA）获取一个SSL证书，我们需要生成一个证书签名请求（Certificate Signing Reqeusts，CRSs）。一个CSR主要包含钥匙对中的公钥，以及其它一些重要的信息。

- 我们可以根据前面生成的私钥生成一个CSR：
  
```gotemplate
openssl req -new -sha256 -key server.key -out server.csr
```

- 执行上面的命令后，需要完成一些信息的填写，主要有：
  
```gotemplate
Country Name (2 letter code) []:CN
State or Province Name (full name) []:xxxx
Locality Name (eg, city) []:xxxx
Organization Name (eg, company) []:xxxx
Organizational Unit Name (eg, section) []:xxxx
Common Name (eg, fully qualified host name) []:xxxx
Email Address []:xxxx@qq.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:123456
```
填写完这些信息后，就会生成一个证书签名请求（server.csr）。

#### 生成证书

如果想使用一个SSL证书来对通信进行加密，但是不需要使用CA签字的证书，那么我们可以生成一个自签名的证书。

使用前面生成的私钥（server.key）以及证书签名请求（server.csr），我们可以生成一个自签名的证书：

```gotemplate
openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 3650
```
选项-x509指定req来生成一个自签名的证书。-days 3650指定了证书的有效期是3650天。-signkey指定了私钥，而-in指定了证书签名请求。

这样，就能生成一个自签名的证书（server.crt）。

#### 注意Common Name


有一个需要注意的问题是，我们在生成服务器时绑定的地主必须和在生成CSR时候填写的Common Name信息一致，不然会出现下面的错误：
```gotemplate
2018/09/25 17:04:02 cound not compute: rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: authentication handshake failed: x509: certificate is valid for example.com, not localhost"
```
就是说，我们在生成CSR是填写的信息是example.com，但是使用的时候却是localhost，那么这是认证不通过的。因此，我们在生成CSR的时候，Common Name信息需要填成localhost。


和服务器端类似，我们使用credentials.NewClientTLSFromFile()函数创建一个credentials对象creds，并用这个对象创建连接。这个函数也是可变数量参数的函数，可以用来指定多个参数，用来指导建立连接时的行为。
注意，在创建creds的时候我们没有用到私钥（server.key），从名字可知，这个私钥是服务器的。
