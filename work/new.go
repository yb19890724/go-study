// Package new generates micro service templates
package new

import "C"
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
	"github.com/spf13/viper"
	
	"github.com/micro/cli"
	tmpl "github.com/micro/micro/internal/template"
	"github.com/micro/micro/internal/usage"
	"github.com/xlab/treeprint"
	
	//"/micro/internal/db2struct"
	"github.com/micro/micro/internal/db2struct"
	
	"github.com/huandu/xstrings"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ungerik/go-dry"
	"gopkg.in/yaml.v2"
	"net/url"
)

type Database struct {
	//root:123456@tcp(10.10.200.12:10191)/kkgoo_3306_kkgoo?tablename=kk_user_settings&prefix=kk_&charset=utf8mb4&parseTime=True&loc=Local
	Dsn string
	//kkgoo_3306_kkgoo
	DbName string
	//kk_user_settings
	Tablename string
	//kk_user_settings
	RealTablename string
	//kk_
	Prefix string
	//_128
	Suffix string
	
	ShardKey string
	ShardType string
	ShardKeyType string
	SpliTableCount int
	
	//userSettings
	Model string
	//
	Key      map[string][]string
	KeyProto map[string]string
	Proto    string
	Column   *map[string]map[string]string
	Rds      string
}

type config struct {
	Database map[string]Database
	//10.10.200.12:10192
	RedisDsn string
	// foo
	Alias string
	// micro new example -type
	Command string
	// go.micro
	Namespace string
	// api, srv, web, fnc
	Type string
	// go.micro.api.foo
	FQDNAPI string
	// go.micro.srv.foo
	FQDNSRV string
	// go.micro.cli.foo
	FQDNCLI string
	// go.micro.web.foo
	FQDNWEB string
	// github.com/micro/foo
	Dir string
	// $GOPATH/src/github.com/micro/foo
	GoDir string
	// $GOPATH
	GoPath string
	// Files
	Files []file
	// Comments
	Comments []string
	// Plugins registry=etcd:broker=nats
	Plugins []string
	//
	Handlers   []string
	Subscriber []string
	Extra      map[string]interface{}
	Srv      []string
}

type file struct {
	Path  string
	Tmpl  string
	Extra map[string]interface{}
	Rewrite bool
}

func write(c config, file, tmpl string) error {
	fn := template.FuncMap{
		"title":            strings.Title,
		"replace":          strings.Replace,
		"ToLower":          strings.ToLower,
		"toCamelCase":      xstrings.ToCamelCase,
		"ToKebabCase":      xstrings.ToKebabCase,
		"FirstRuneToLower": xstrings.FirstRuneToLower,
	}
	
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	
	t, err := template.New("f").Funcs(fn).Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, c)
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.GoDir); !os.IsNotExist(err) {
		//return fmt.Errorf("%s already exists", c.GoDir)
	}
	
	// create usage report
	u := usage.New("new")
	// a single request/service
	u.Metrics.Count["requests"] = uint64(1)
	u.Metrics.Count["services"] = uint64(1)
	// send report
	go usage.Report(u)
	
	// just wait
	<-time.After(time.Millisecond * 250)
	
	fmt.Printf("Creating service %s in %s\n\n", c.Namespace, c.GoDir)
	
	t := treeprint.New()
	
	nodes := map[string]treeprint.Tree{}
	nodes[c.GoDir] = t
	
	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.GoDir, file.Path)
		dir := filepath.Dir(f)
		
		b, ok := nodes[dir]
		if !ok {
			d, _ := filepath.Rel(c.GoDir, dir)
			b = t.AddBranch(d)
			nodes[dir] = b
		}
		
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		
		p := filepath.Base(f)
		
		b.AddNode(p)
		c.Extra = file.Extra
		_, err := os.Stat(f)
		if !file.Rewrite && (err == nil || os.IsExist(err)) {
			continue
		}
		if err := write(c, f, file.Tmpl); err != nil {
			fmt.Printf("Creating %s error\n", f)
			return err
		}
	}
	
	// print tree
	fmt.Println(t.String())
	
	for _, comment := range c.Comments {
		fmt.Println(comment)
	}
	
	// just wait
	<-time.After(time.Millisecond * 250)
	
	return nil
}

func getProtoService(proto string, serviceName string) string {
	re := regexp.MustCompile("service "+ serviceName +" {[\\S\\s]+?}\n")
	proto = re.FindString(proto)
	proto = strings.TrimSpace(proto)
	re = regexp.MustCompile("{[\\S\\s]*}$")
	proto = re.FindString(proto)
	proto = strings.TrimSpace(proto)
	proto = strings.Replace(proto, "{\n}", "", 1)
	proto = strings.Replace(proto, "{\n", "", 1)
	proto = strings.Replace(proto, "\n}", "", 1)
	return proto
}

func getProtoMessage(proto string) string {
	re := regexp.MustCompile("message [^{]* {[\\S\\s]*}\n")
	all := re.FindAllString(proto, -1)
	protoMessage := strings.Join(all, "\n")
	return protoMessage
}


// key [id uid]
// LikeUserAvatarIdAndUidKey
func filterProtoByKey(proto string, key []string, message string) string {
	proto = strings.TrimSpace(proto)
	re := regexp.MustCompile("{[^}]*}$")
	proto = re.FindString(proto)
	proto = strings.Replace(proto, "{", "", 1)
	proto = strings.Replace(proto, "}", "", 1)
	proto = strings.TrimSpace(proto)
	
	arrProto := strings.Split(proto, ";")
	arrProto = dry.StringMap(strings.TrimSpace, arrProto)
	protoMap := make(map[string]string)
	for _, protoKey := range arrProto {
		list := strings.Split(protoKey, " = ")
		listTypeKey := strings.Split(list[0], " ")
		if len(listTypeKey) == 2 {
			protoMap[listTypeKey[1]] = listTypeKey[0]
		}
	}
	
	var protoKeyString string
	var order int
	order = 1
	protoKeyString = "message " + message + " {"
	for _, k := range key {
		if _, ok := protoMap[k]; ok {
			protoKeyString += "\n\t" + protoMap[k] + " " + k + " = " + strconv.Itoa(order) + ";"
			order += 1
		}
	}
	protoKeyString += "\n}"
	return protoKeyString
}

func getProtoTypeByKey(proto string, key string) string {
	proto = strings.TrimSpace(proto)
	re := regexp.MustCompile("{[^}]*}$")
	proto = re.FindString(proto)
	proto = strings.Replace(proto, "{", "", 1)
	proto = strings.Replace(proto, "}", "", 1)
	proto = strings.TrimSpace(proto)
	
	arrProto := strings.Split(proto, ";")
	arrProto = dry.StringMap(strings.TrimSpace, arrProto)
	protoMap := make(map[string]string)
	for _, protoKey := range arrProto {
		list := strings.Split(protoKey, " = ")
		listTypeKey := strings.Split(list[0], " ")
		if len(listTypeKey) == 2 {
			protoMap[listTypeKey[1]] = listTypeKey[0]
		}
	}
	
	if _, ok := protoMap[key]; ok {
		return protoMap[key]
	}
	return ""
}

func run(ctx *cli.Context) {
	
	// 读取设置配置
	namespace := ctx.String("namespace")
	alias := ctx.String("alias")
	dir := ctx.Args().First()
	useGoPath := ctx.Bool("gopath")
	dbdsn := ctx.StringSlice("dbdsn")
	useGoModule := os.Getenv("GO111MODULE")
	
	var plugins []string
	
	
	// 判断服务名称 likeUserAvatar
	if len(dir) == 0 {
		fmt.Println("specify service name")
		return
	}
	
	// go.micro
	if len(namespace) == 0 {
		fmt.Println("namespace not defined")
		return
	}
	
	// set the command
	command := fmt.Sprintf("micro new %s", dir)
	if len(namespace) > 0 {
		command += " --namespace=" + namespace
	}
	if len(alias) > 0 {
		command += " --alias=" + alias
	}
	
	if plugins := ctx.StringSlice("plugin"); len(plugins) > 0 {
		command += " --plugin=" + strings.Join(plugins, ":")
	}
	
	if len(dbdsn) > 0 {
		for _, dsn := range dbdsn {
			command += " --dbdsn=\"" + dsn + "\""
		}
	}
	
	//command :micro new likeUserAvatar --namespace=go.micro --dbdsn="root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?tablename=kk_like_user_avatar_0&prefix=kk_&suffix0&charset=utf8mb4&parseTime=True&loc=Local&shardKey=uid&shardType=mod&spliTableCount=128&rds=10.10.200.12:10058"
	
	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(dir) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return
	}
	
	var goPath string
	var goDir string
	var error error
	
	// only set gopath if told to use it
	
	// 是否使用 gopath 读取系统gopath目录
	if useGoPath {
		goPath = os.Getenv("GOPATH")
		
		// don't know GOPATH, runaway....
		if len(goPath) == 0 {
			fmt.Println("unknown GOPATH")
			return
		}
		
		// attempt to split path if not windows
		if runtime.GOOS == "windows" {
			goPath = strings.Split(goPath, ";")[0]
		} else {
			goPath = strings.Split(goPath, ":")[0]
		}
		
		goDir = filepath.Join(goPath, "src", path.Clean(dir))//Users/nice/go/src/likeUserAvatar
		
	} else {
		goDir = path.Clean(dir)
	}
	
	if len(alias) == 0 {
		// set as last part
		alias = filepath.Base(dir)
	}
	
	alias = strings.Replace(alias, "-", "_", -1)
	
	// go.micro.srv.likeUserAvatar
	// go.micro.web.likeUserAvatar
	// go.micro.api.likeUserAvatar
	// go.micro.cli.likeUserAvatar
	
	fqdnSrv := strings.Join([]string{namespace, "srv", alias}, ".")
	fqdnApi := strings.Join([]string{namespace, "api", alias}, ".")
	fqdnWeb := strings.Join([]string{namespace, "web", alias}, ".")
	fqdnCli := strings.Join([]string{namespace, "cli", alias}, ".")
	
	for _, plugin := range ctx.StringSlice("plugin") { //broker 提供者
		// registry=etcd:broker=nats
		for _, p := range strings.Split(plugin, ":") {
			// registry=etcd
			parts := strings.Split(p, "=")
			if len(parts) < 2 {
				continue
			}
			plugins = append(plugins, path.Join(parts...))
		}
	}
	
	///Users/nice/go/src/likeUserAvatar/proto
	os.RemoveAll(filepath.Join(goDir, "proto"))//删除目录以及其子目录和文件，如果path不存在的话，返回nil
	
	if _, err := os.Stat(filepath.Join(goDir, "proto")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(goDir, "proto"), 0755); err != nil {
			fmt.Println("mkdir proto error")
			return
		}
	}
	
	database := make(map[string]Database)
	
	if len(dbdsn) > 0 {
		storageDir := filepath.Join(goDir, "srv", "storage")
		// 创建 storage 文件夹
		if _, err := os.Stat(storageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(storageDir, 0755); err != nil {
				fmt.Println("mkdir storage error")
				return
			}
		}
		for _, dsn := range dbdsn {
			db := new(Database)
			db.Key = make(map[string][]string)
			db.KeyProto = make(map[string]string)
			u, err := url.Parse(dsn)
			
			if err != nil {
				fmt.Println("dbdsn parse error")
				return
			}
			db.DbName = path.Base(u.Opaque)
			if "" == db.DbName {
				fmt.Println("dbname parse error")
				return
			}
			
			//query, err := url.ParseQuery(u.RawQuery)
			query := u.Query()
			if err != nil {
				fmt.Println("dbdsn parse error")
				return
			}
			
			if v, ok := query["tablename"]; ok {
				db.RealTablename = v[0]
				db.Tablename = v[0]
				query.Del("tablename")
				
				if v2, ok := query["prefix"]; ok {
					re, _ := regexp.Compile("^" + v2[0])
					db.Prefix = v2[0]
					db.Tablename = re.ReplaceAllString(db.Tablename, "")
					query.Del("prefix")
				} else {
					db.Prefix = ""
				}
				
				if v2, ok := query["suffix"]; ok {
					re, _ := regexp.Compile(v2[0] + "$")
					db.Suffix = v2[0]
					db.Tablename = re.ReplaceAllString(db.Tablename, "")
					query.Del("suffix")
				} else {
					db.Suffix = ""
				}
				db.Model = xstrings.ToCamelCase(db.Tablename)
				
				if v2, ok := query["shardKey"]; ok {
					db.ShardKey = v2[0]
					query.Del("shardKey")
				}
				
				if v2, ok := query["shardType"]; ok {
					db.ShardType = v2[0]
					query.Del("shardType")
				}
				
				if v2, ok := query["spliTableCount"]; ok {
					db.SpliTableCount,_ = strconv.Atoi(v2[0])
					query.Del("spliTableCount")
				}
			} else {
				fmt.Println("need tablename to build svr")
				return
			}
			
			if v, ok := query["rds"]; ok {
				db.Rds = v[0]
				query.Del("rds")
			} else {
				fmt.Println("need rds to build svr")
				return
			}
			u.RawQuery = query.Encode()
			db.Dsn = u.String()
			
			
			// columnDataTypes : map 字段设置
			// columnKeyTypes : &map[PRIMARY:[id] uid_avatarTime_likeUid:[uid avatar_time like_uid]]
			// createTable :建表语句
			
			columnDataTypes, columnKeyTypes, createTable, err := db2struct.GetColumnsFromMysqlTable(db.Dsn, db.DbName, db.RealTablename)
			
			if err != nil {
				fmt.Println("db 2 struct error", err)
				return
			}
			
			db.Column = columnDataTypes
			
			f, err := os.Create(filepath.Join(storageDir, db.Model+".sql"))
			f.Write([]byte(createTable))
			f.Close()
			
			// 创建sql文件  mysql-protobuf       /Users/nice/go/src/likeUserAvatar/srv/storage/LikeUserAvatar.sql
			protoCmd := exec.Command("mysql-protobuf", filepath.Join(storageDir, db.Model+".sql"))
			protoOut, err := protoCmd.CombinedOutput() //运行命令
			if err != nil {
				fmt.Println("db 2 proto error", err, protoOut)
				return
			}
			db.Proto = strings.Replace(string(protoOut), "message "+strings.Title(db.RealTablename), "message "+db.Model, 1)
			/**
			    syntax = "proto3";
				message LikeUserAvatar {
				  uint32 id = 1;
				  uint32 uid = 2;
				  uint32 avatar_time = 3;
				  uint32 like_uid = 4;
				  uint32 like_time = 5;
				}
			*/
			db.Proto = strings.Replace(db.Proto, "syntax = \"proto3\";", "", 1)
			
			
			for _, v := range *columnKeyTypes { // 循环 proto文件内容  内容
				if "" != db.ShardKey && false == dry.StringInSlice(db.ShardKey, v) {
					v = append(v, db.ShardKey)
				}
				db.Key[strings.Join(dry.StringMap(dry.StringToUpperCamelCase, v), "And")] = v
				db.KeyProto[strings.Join(dry.StringMap(dry.StringToUpperCamelCase, v), "And")] = filterProtoByKey(db.Proto, v, db.Model+strings.Join(dry.StringMap(dry.StringToUpperCamelCase, v), "And")+"Key")
			}
			/**
				map[IdAndUid:message LikeUserAvatarIdAndUidKey {
			        uint32 id = 1;
			        uint32 uid = 2;
			} UidAndAvatarTimeAndLikeUid:message LikeUserAvatarUidAndAvatarTimeAndLikeUidKey {
			        uint32 uid = 1;
			        uint32 avatar_time = 2;
			        uint32 like_uid = 3;
			} ShardKey:message ShardKey {
			        uint32 uid = 1;
			}]
			*/
			db.KeyProto["ShardKey"] = filterProtoByKey(db.Proto, []string{db.ShardKey}, "ShardKey")
			db.ShardKeyType = getProtoTypeByKey(db.Proto, db.ShardKey) // uint32
			
			// model struct template string
			struc, err := db2struct.Generate(*columnDataTypes, db.Prefix + db.Tablename, db.Model, "storage", false, true, false, db.ShardType, db.ShardKey, db.ShardKeyType, db.SpliTableCount)
			structPath := filepath.Join(storageDir, db.Model+".go")
			
			f, err = os.Create(structPath)
			f.Write([]byte(string(struc)))
			f.Close()
			
			database[db.Model] = *db
		}
	}
	
	var c config
	
	// create srv config
	c = config{
		Alias:     alias,
		Command:   command,
		Namespace: namespace,
		FQDNAPI:   fqdnApi,
		FQDNSRV:   fqdnSrv,
		FQDNCLI:   fqdnCli,
		FQDNWEB:   fqdnWeb,
		Dir:       dir,
		GoDir:     goDir,
		GoPath:    goPath,
		Plugins:   plugins,
		Database:  database,
		Files: []file{
			{"srv/main.go", tmpl.MainSRV, nil, true},
			{".gitignore", tmpl.GitIgnore, nil, true},
			{"connect/log.go", tmpl.ConnectLogSRV, nil, true},
			{"connect/config.go", tmpl.ConnectConfigSRV, nil, true},
			{"k8s/.helmignore", tmpl.K8sHelmignore, nil, true},
			{"k8s/Chart.yaml", tmpl.K8sChart, nil, true},
			{"k8s/templates/service.yaml", tmpl.K8sService, nil, true},
			{"k8s/templates/NOTES.txt", tmpl.K8sNotes, nil, true},
			{"k8s/templates/_helpers.tpl", tmpl.K8sHelpers, nil, true},
			{"k8s/templates/deployment.yaml", tmpl.K8sDeployment, nil, true},
			{"k8s/templates/ingress.yaml", tmpl.K8sIngress, nil, true},
			{"k8s/templates/tests/test-connection.yaml", tmpl.K8sTestConnection, nil, true},
			{"plugin.go", tmpl.Plugin, nil, true},
			{"srv/Dockerfile", tmpl.DockerSRV, nil, true},
			{"deploy.sh", tmpl.Deploy, nil, true},
			{"helper/helper.go", tmpl.HelperStruct, nil, true},
			{"helper/pprof.go", tmpl.HelperPprof, nil, true},
			{"helper/timer.go", tmpl.HelperTimer, nil, true},
			{"helper/path.go", tmpl.HelperPath, nil, true},
			{"Makefile", tmpl.Makefile, nil, true},
			{"README.md", tmpl.Readme, nil, true},
		},
		Comments: []string{
			"\ndownload protobuf for micro:\n",
			"brew install protobuf",
			"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
			"go get -u github.com/micro/protoc-gen-micro",
			"\ncompile the proto file " + alias + ".proto:\n",
			"cd " + goDir,
			"protoc --proto_path=. --go_out=. --micro_out=. proto/" + alias + ".proto\n",
		},
		Srv: []string{
		},
	}
	
	// /Users/nice/go/src/likeUserAvatar/srv.yaml
	srvConf := filepath.Join(c.GoDir, "srv.yaml")
	
	// strSrvConf : []
	// err        : open /Users/nice/go/src/likeUserAvatar/srv.yaml: no such file or directory
	strSrvConf, err := ioutil.ReadFile(srvConf)
	
	// 因为错误demo 没有走里面
	if err == nil {
		viperObj := viper.New()
		viperObj.SetConfigType("yaml")
		err = viperObj.ReadConfig(bytes.NewBuffer(strSrvConf))
		if err != nil {
			fmt.Println("srv.yaml error", err)
			return
		}
		allSrv := viperObj.AllSettings()
		defer os.RemoveAll(filepath.Join(c.GoDir, "build"))
		for k,_ := range allSrv {
			os.RemoveAll(filepath.Join(c.GoDir, "build"))
			gitCmd := exec.Command("git", "clone", viperObj.GetString(k+".git"), filepath.Join(c.GoDir, "build"), "--branch="+viperObj.GetString(k+".branch"))
			gitOut, err := gitCmd.CombinedOutput()
			if err != nil {
				fmt.Println("git clone srv error", viperObj.GetString(k+".git"), err, string(gitOut))
				return
			}
			protoDir, _ := filepath.Glob(filepath.Join(c.GoDir, "build", "srv", "proto", "*"))
			for _, v := range protoDir {
				os.Rename(filepath.Join(c.GoDir, "build", "srv", "proto", path.Base(v)), filepath.Join(c.GoDir, "proto", path.Base(v)))
				c.Srv = append(c.Srv, path.Base(v))
			}
		}
	}
	
	// 生成 mysql 和 redis go文件
	var strRedisYaml []byte
	var strDatabaseYaml []byte
	if len(dbdsn) > 0 {
		c.Files = append(c.Files, file{"connect/mysql.go", tmpl.ConnectMysqlSRV, nil, true})
		c.Files = append(c.Files, file{"connect/redis.go", tmpl.ConnectRedisSRV, nil, true})
		databaseYaml := make(map[string]interface{})
		rdsYaml := make(map[string]interface{})
		keys := make(map[string][]string)
		
		/*
				database: map[LikeUserAvatar:{root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True nice_23313_user_avatar like_user_avatar kk_like_user_avatar_0 kk_ _0 uid mod uint32 128 LikeUserAvatar map[UidAndAvatarTimeAndLikeUid:[uid avatar_time like_uid] IdAndUid:[id uid]] map[UidAndAvatarTimeAndLikeUid:message LikeUserAvatarUidAndAvatarTimeAndLikeUidKey {
		                uint32 uid = 1;
		                uint32 avatar_time = 2;
		                uint32 like_uid = 3;
		        } IdAndUid:message LikeUserAvatarIdAndUidKey {
		                uint32 id = 1;
		                uint32 uid = 2;
		        } ShardKey:message ShardKey {
		                uint32 uid = 1;
		        }] */
		
		for _, v := range database {
			/*
					        for range v : {root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True nice_23313_user_avatar like_user_avatar kk_like_user_avatar_0 kk_ _0 uid mod uint32 128 LikeUserAvatar map[UidAndAvatarTimeAndLikeUid:[uid avatar_time like_uid] IdAndUid:[id uid]] map[UidAndAvatarTimeAndLikeUid:message LikeUserAvatarUidAndAvatarTimeAndLikeUidKey {
			                        uint32 uid = 1;
			                        uint32 avatar_time = 2;
			                        uint32 like_uid = 3;
			                } IdAndUid:message LikeUserAvatarIdAndUidKey {
			                        uint32 id = 1;
			                        uint32 uid = 2;
			                } ShardKey:message ShardKey {
			                        uint32 uid = 1;
			                }]
			
			*/
			item := make(map[string]interface{})
			item = map[string]interface{}{"dsn": v.Dsn, "max_idle_conns": 2, "max_open_conns": 4, "conn_max_lifetime": 300}
			databaseYaml[v.Model] = map[string]interface{}{"master": item, "slave": item}
			rdsYaml[v.Model] = map[string]interface{}{"addr": v.Rds, "maxRetries": 1, "dialTimeout": "100ms", "poolSize": 20, "readTimeout": "100ms", "writeTimeout": "100ms", "minIdleConns": 10, "maxConnAge": "5min"}
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, "Init.go"), tmpl.InitSRV, map[string]interface{}{"model": v.Model, "srv":c.Srv}, true})
			
			//合并两个proto文件
			protoCustomizeService := ""
			protoCustomizeMessage := ""
			protoCustomizeFile := filepath.Join(c.GoDir, "srv", "proto", v.Model, v.Model+".customize.proto")
			//protoCustomizeFile /Users/nice/go/src/likeUserAvatar/srv/proto/LikeUserAvatar/LikeUserAvatar.customize.proto
			protoCustomize, err := ioutil.ReadFile(protoCustomizeFile)
			if err == nil {
				protoCustomizeService = getProtoService(string(protoCustomize), v.Model+"SRV")
				protoCustomizeMessage = getProtoMessage(string(protoCustomize))
			}
			
			// 加入创建文件
			c.Files = append(c.Files, file{filepath.Join("srv", "proto", v.Model, v.Model+".proto"), tmpl.ProtoSRV, map[string]interface{}{"model": v.Model, "protoCustomizeService": protoCustomizeService, "protoCustomizeMessage":protoCustomizeMessage}, true})
			c.Files = append(c.Files, file{filepath.Join("srv", "proto", v.Model, v.Model+".customize.proto"), tmpl.ProtoCustomizeSRV, map[string]interface{}{"model": v.Model}, false})
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, "Create"+v.Model+".go"), tmpl.CreateSRV, map[string]interface{}{"model": v.Model, "func": "Create" + v.Model}, true})
			
			var fun string
			var relatedfunc string
			for keyName, _ := range v.Key {
				//v.key: map[UidAndAvatarTimeAndLikeUid:[uid avatar_time like_uid] IdAndUid:[id uid]]
				//keyName: IdAndUid
				keys[v.Model] = append(keys[v.Model], keyName)
				fun = "Delete" + v.Model + "By" + keyName
				//fun: DeleteLikeUserAvatarByIdAndUid`
				c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.DeleteSRV, map[string]interface{}{"model": v.Model, "func": fun, "keyName": keyName}, true})
				fun = "Get" + v.Model + "By" + keyName
				c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.GetKeySRV, map[string]interface{}{"model": v.Model, "func": fun, "keyName": keyName}, true})
				
				relatedfunc = fun
				fun = "Mget" + v.Model + "By" + keyName
				c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.MgetKeySRV, map[string]interface{}{"model": v.Model, "func": fun, "relatedfunc": relatedfunc, "keyName": keyName}, true})
				fun = "Update" + v.Model + "By" + keyName
				c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.UpdateKeySRV, map[string]interface{}{"model": v.Model, "func": fun, "keyName": keyName}, true})
				fun = "Replace" + v.Model + "By" + keyName
				c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.ReplaceKeySRV, map[string]interface{}{"model": v.Model, "func": fun, "keyName": keyName}, true})
			}
			fun = "Get" + v.Model + "ListByWhere"
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.GetWhereSRV, map[string]interface{}{"model": v.Model, "func": fun}, true})
			fun = "Get" + v.Model + "CountByWhere"
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.GetWhereCountSRV, map[string]interface{}{"model": v.Model, "func": fun}, true})
			fun = "Get" + v.Model + "ListByAssoc"
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.GetAssocSRV, map[string]interface{}{"model": v.Model, "func": fun}, true})
			fun = "Get" + v.Model + "CountByAssoc"
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", v.Model, fun+".go"), tmpl.GetAssocCountSRV, map[string]interface{}{"model": v.Model, "func": fun}, true})
		}
		
		//map[LikeUserAvatar:map[master:map[max_idle_conns:2 max_open_conns:4 conn_max_lifetime:300 dsn:root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True] slave:map[dsn:root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True max_idle_conns:2 max_open_conns:4 conn_max_lifetime:300]]]
		
		strDatabaseYaml, error = yaml.Marshal(&databaseYaml)
		if error != nil {
			fmt.Println("convert to yaml fail", error)
			return
		}
		// "yaml":LikeUserAvatar:
	master:
	conn_max_lifetime: 300
	dsn: root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True
	max_idle_conns: 2
	max_open_conns: 4
	slave:
	conn_max_lifetime: 300
	dsn: root:123456@tcp(10.10.200.12:10055)/nice_23313_user_avatar?charset=utf8mb4&loc=Local&parseTime=True
	max_idle_conns: 2
	max_open_conns: 4
		
		c.Files = append(c.Files, file{"config/database.yaml", tmpl.Yaml, map[string]interface{}{"yaml": string(strDatabaseYaml)}, true})
		strRedisYaml, error = yaml.Marshal(&rdsYaml)
		if error != nil {
			fmt.Println("convert to yaml fail", error)
			return
		}
		c.Files = append(c.Files, file{"config/redis.yaml", tmpl.Yaml, map[string]interface{}{"yaml": string(strRedisYaml)}, true})
	} else {
		//合并两个proto文件
		protoCustomizeService := ""
		protoCustomizeMessage := ""
		protoCustomizeFile := filepath.Join(c.GoDir, "srv", "proto", xstrings.ToCamelCase(alias), xstrings.ToCamelCase(alias)+".customize.proto")
		protoCustomize, err := ioutil.ReadFile(protoCustomizeFile)
		if err == nil {
			protoCustomizeService = getProtoService(string(protoCustomize), xstrings.ToCamelCase(alias)+"SRV")
			protoCustomizeMessage = getProtoMessage(string(protoCustomize))
		}
		
		c.Files = append(c.Files, file{filepath.Join("srv", "proto", xstrings.ToCamelCase(alias), xstrings.ToCamelCase(alias)+".proto"), tmpl.ProtoLogicSRV, map[string]interface{}{"model": xstrings.ToCamelCase(alias), "protoCustomizeService": protoCustomizeService, "protoCustomizeMessage":protoCustomizeMessage}, true})
		c.Files = append(c.Files, file{filepath.Join("srv", "proto", xstrings.ToCamelCase(alias), xstrings.ToCamelCase(alias)+".customize.proto"), tmpl.ProtoCustomizeSRV, map[string]interface{}{"model": xstrings.ToCamelCase(alias)}, false})
	}
	
	logYaml := make(map[string]interface{})
	logYaml["dirpath"] = filepath.Join("var", "log")
	logYaml["level"] = "trace"
	strLogYaml, err := yaml.Marshal(&logYaml)
	if err != nil {
		fmt.Println("convert to yaml fail", err)
		return
	}
	
	//"yaml": string(strLogYaml) 写入文件的字符串
	/*dirpath: var/log
	  level: trace*/
	c.Files = append(c.Files, file{"config/log.yaml", tmpl.Yaml, map[string]interface{}{"yaml": string(strLogYaml)}, true})
	
	configMap := make(map[string]map[string]interface{})
	configMap["configmap"] = make(map[string]interface{})
	if len(dbdsn) > 0 {
		configMap["configmap"]["redis.yaml"] = string(strRedisYaml)
		configMap["configmap"]["database.yaml"] = string(strDatabaseYaml)
	}
	configMap["configmap"]["log.yaml"] = string(strLogYaml)
	strConfigMapYaml, err := yaml.Marshal(&configMap)
	
	c.Files = append(c.Files, file{"k8s/values.yaml", tmpl.K8sValues, map[string]interface{}{"configMapYaml": string(strConfigMapYaml)}, false})
	c.Files = append(c.Files, file{"k8s/templates/configmap.yaml", tmpl.K8sConfigmap, nil, true})
	
	// set gomodule
	if useGoModule == "on" || useGoModule == "auto" {
		c.Files = append(c.Files, file{"go.mod", tmpl.Module, nil, true})
	}
	
	if err := create(c); err != nil {
		fmt.Println(err)
		return
	}
	
	// 生成微服务
	protoDir, _ := filepath.Glob(filepath.Join(c.GoDir, "srv", "proto", "*", "*.proto"))
	for _, protoFile := range protoDir {
		if -1 != strings.Index(protoFile, ".customize.proto") {
			continue
		}
		protoCmd := exec.Command("protoc", "--proto_path="+c.GoDir, "--go_out="+c.GoDir, "--micro_out="+c.GoDir, protoFile)
		protoOut, err := protoCmd.CombinedOutput()
		if err != nil {
			fmt.Println("proto to go error", protoFile, err, string(protoOut))
			return
		}
	}
	
	pbDir, _ := filepath.Glob(filepath.Join(c.GoDir, "srv", "proto", "*", "*.pb.go"))
	for _, pbFile := range pbDir {
		sedCmd := exec.Command("sed", "-i", "", "s/,omitempty//g", pbFile)
		sedOut, err := sedCmd.CombinedOutput()
		if err != nil {
			fmt.Println("sed error", err, string(sedOut))
			return
		}
	}
	
	//protoDir [/Users/nice/go/src/likeUserAvatar/srv/proto/LikeUserAvatar/LikeUserAvatar.micro.go]
	protoDir, _ = filepath.Glob(filepath.Join(c.GoDir, "srv", "proto", "*", "*.micro.go"))
	
	for _, v := range protoDir {
		f, err := os.Open(v)
		if err != nil {
			fmt.Println("read proto error", v, err)
			return
		}
		
		//protoReadFile ：LikeUserAvatar.micro.go文件字符串内容
		protoReadFile, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("read proto error", v, err)
			return
		}
		
		reg := regexp.MustCompile(`type ([\w]+)Service interface`)
		funs := reg.FindAllStringSubmatch(string(protoReadFile), -1)
		model := strings.Replace(funs[0][1], "SRV", "", -1)
		
		reg = regexp.MustCompile(`([\w]+)\(ctx context.Context, in`)
		funs = reg.FindAllStringSubmatch(string(protoReadFile), -1)
		
		funcMap := make(map[string]interface{})
		for _, v2 := range funs {
			fun := v2[1]
			funcMap[v2[1]] = v2[1]
			_, err := os.Stat(filepath.Join("srv", "handler", model, fun+".go"))
			if (err == nil || os.IsExist(err)) {
				continue
			}
			c.Files = append(c.Files, file{filepath.Join("srv", "handler", model, fun+".go"), tmpl.HandlerSRV, map[string]interface{}{"model": model, "func": fun}, false})
		}
		c.Files = append(c.Files, file{filepath.Join("srv", "proto", model, model+".pool.go"), tmpl.ProtoPool, map[string]interface{}{"model": model, "funcs": funcMap}, true})
		c.Files = append(c.Files, file{filepath.Join("srv", "proto", model, "example_test.go"), tmpl.ProtoPoolTest, map[string]interface{}{"model": model, "funcs": funcMap}, true})
	}
	
	if err := create(c); err != nil {
		fmt.Println(err)
		return
	}
	
	protoCmd := exec.Command("gofmt", "-s", "-w", goDir)
	protoOut, err := protoCmd.CombinedOutput()
	if err != nil {
		fmt.Println("gofmt error", err, string(protoOut))
		return
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "new",
			Usage: "Create a service template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Usage: "Namespace for the service e.g com.example",
					Value: "go.micro",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "Type of service e.g api, fnc, srv, web",
					Value: "srv",
				},
				cli.StringFlag{
					Name:  "fqdn",
					Usage: "FQDN of service e.g com.example.srv.service (defaults to namespace.type.alias)",
				},
				cli.StringFlag{
					Name:  "alias",
					Usage: "Alias is the short name used as part of combined name if specified",
				},
				cli.StringSliceFlag{
					Name:  "plugin",
					Usage: "Specify plugins e.g --plugin=registry=etcd:broker=nats or use flag multiple times",
				},
				cli.StringSliceFlag{
					Name:  "dbdsn",
					Usage: "database dsn of service e.g --dbdsn=root:123456@tcp(10.10.200.12:10191)/kkgoo_3306_kkgoo?tablename=kk_user_settings&charset=utf8mb4&parseTime=True&loc=Local and use flag multiple times",
				},
				
				cli.BoolTFlag{
					Name:  "gopath",
					Usage: "Create the service in the gopath. Defaults to true.",
				},
			},
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
