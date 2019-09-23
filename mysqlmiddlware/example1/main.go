package main


import (
	"fmt"
	"github.com/siddontang/go-mysql/client"
	_ "github.com/siddontang/go-mysql/driver"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/xwb1989/sqlparser"
	"log"
)

type MyHandle struct {
	server.Handler
	Conn *client.Conn
}

// 创建
func NewMyHandle() MyHandle{
	
	conn, err := client.Connect("127.0.0.1:3306", "root", "root", "test")
	if err!= nil {
		log.Fatal(err)
	}
	return MyHandle{Conn:conn}
}

// 代理执行sql
func (mh MyHandle) HandleQuery(query string) (*mysql.Result, error) {
	
	result ,err:=mh.Conn.Execute(query)
	if err!=nil {
		return nil, fmt.Errorf("error:",query)
	}
	return result, nil
}

//


// 监听tcp链接 模拟启动一个mysql服务
// 服务器包提供了一个框架来实现一个简单的MySQL服务器，它可以处理来自MySQL客户端的数据包。
// 您可以使用它来构建自己的MySQL代理。服务器连接与MySQL 5.5,5.6,5.7和8.0版本兼容，因此大多数MySQL客户端应该能够无需修改即可连接到服务器。
func main()  {
	
	// l, _ := net.Listen("tcp", "0.0.0.0:4000")
	//
	// for   { // 循环是为了接送多个链接
	//
	// 	c, _ := l.Accept()// 有客户端成功链接
	//
	// 	go func() {
	//
	// 		// Create a connection with user root and an empty password.
	// 		// You can use your own handler to handle command here.
	// 		//conn, _ := server.NewConn(c, "root", "123", server.EmptyHandler{})
	// 		conn, _ := server.NewConn(c, "root", "123", NewMyHandle())
	//
	// 		for {
	// 			conn.HandleCommand() // 接受命令发送给 EmptyHandler
	// 		}
	// 	}()
	//
	// }
	sqlParse()
}

func sqlParse()  {
	sql :="select *,id,name from users as a "

	stmt,err:=sqlparser.Parse(sql)
	
	if err != nil {
		log.Fatal(err)
	}
	
	switch stmt:=stmt.(type) {
	
	case *sqlparser.Select:
		buff:=sqlparser.NewTrackedBuffer(nil)
		stmt.SelectExprs.Format(buff)
		fmt.Println(buff.String())
		
		// for _,node:= range stmt.From{
		// 	getTable := node.(*sqlparser.AliasedTableExpr)
		// 	fmt.Println(getTable.As.String())
		// 	fmt.Println(getTable.Expr.(sqlparser.TableName).Name)
		// }
	case *sqlparser.Insert:
	
	case *sqlparser.Delete:
	
	}
}