package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	v1 "grpcpro/api/service/v1"
	"grpcpro/server"
)

type Config struct {
	GRPCPort string
	DataStoreDBHost string
	DataStoreDBUser string
	DataStoreDBPassword string
	DataStoreDBSchema string
}

func RunServer()error{
	ctx:= context.Background()
	var cfg Config
	flag.StringVar(&cfg.GRPCPort,"grpc-port", "9090","gRPC port to bind")
	flag.StringVar(&cfg.DataStoreDBHost,"db-host","127.0.0.1:3306","db-host")
	flag.StringVar(&cfg.DataStoreDBUser,"db-user","root","db-user")
	flag.StringVar(&cfg.DataStoreDBPassword,"db-password","sjy1999","db-password")
	flag.StringVar(&cfg.DataStoreDBSchema,"db-schema","todotask","db-schema")
	flag.Parse()

	if len(cfg.GRPCPort)==0 {
		return fmt.Errorf("invalid TCP port for gRPC server : %s",cfg.GRPCPort)

	}

	param:="parseTime=true"
	dsn:=fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",cfg.DataStoreDBUser,cfg.DataStoreDBPassword,cfg.DataStoreDBHost,cfg.DataStoreDBSchema,param)
	db,err:=sql.Open("mysql",dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v",err)
	}
	defer db.Close()
	v1API:=v1.NewToDoServiceServer(db)
	return server.RunServer(ctx,v1API,cfg.GRPCPort)


}
