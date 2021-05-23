package main

import (
	"context"
	"flag"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	v1 "grpcpro/api/proto/v1"
	"log"
	"time"
)

const apiVersion = "v1" 

func main() {
	address:=flag.String("server","127.0.0.1:9090","gRPC server in format host:post")
	flag.Parse()

	conn,err:= grpc.Dial(*address,grpc.WithInsecure())
	if err != nil {
		log.Fatal("服务端连不上",err)
	}
	defer conn.Close()

	c:=v1.NewToDoServiceClient(conn)  
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	
	t:= time.Now().In(time.UTC)
	reminder ,_:=ptypes.TimestampProto(t)
	pfx:=t.Format(time.RFC3339Nano)
	req1:= v1.CreateRequest{
		Api:  apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title("+pfx+")",
			Description: "description("+pfx+")",
			Reminder:   reminder,
		},
	}
	res1,err:= c.Create(ctx,&req1)
	if err != nil {
		log.Fatalf("创建失败%v",err)
	}
	log.Printf("Create result%v",res1)
	id := res1.Id
	
	req2:=v1.ReadRequest{Api: apiVersion,Id: id}
	res2,err:= c.Read(ctx,&req2)
	if err != nil {
		log.Fatalf("Read failed %v ",err)
	}
	log.Printf("Read result %v",res2)
	
	req3:=v1.UpdateRequest{
		Api:  apiVersion,
		ToDo: &v1.ToDo{
			Id:          res2.ToDo.Id,
			Title:       res2.ToDo.Title,
			Description: res2.ToDo.Description + "updated",
			Reminder:    res2.ToDo.Reminder,
		},
	}
	res3,err := c.Update(ctx,&req3)
	if err != nil {
		log.Fatalf("update failed %v",err)
	}
	log.Printf("update result %v",res3)

	req4:= v1.ReadAllRequest{Api: apiVersion}
	res4,err :=c.ReadAll(ctx,&req4)
	if err != nil {
		log.Fatalf("readall failed %v",err)
	}
	log.Printf("readall result %v",res4)

	req5:=v1.DeleteRequest{
		Api: apiVersion,
		Id:  id,
	}
	res5,err:= c.Delete(ctx,&req5)
	if err != nil {
		log.Fatalf("delete failed %v",err)
	}
	log.Printf("delete result %v",res5)


}
