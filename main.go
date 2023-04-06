package main

import (
	"fmt"

	ros_proto "github.com/HWZen/go_backend/pkg/protobuf"

	ros_hybrid_go "github.com/HWZen/ros_hybrid_go/src"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)


func main(){
	fmt.Println("Hello world!")

	node, err := ros_hybrid_go.NewNode("localhost:5150", "go_backend")
	if err != nil {
		panic(err)
	}
	defer node.Shutdown()

	go node.Run()

	var stats = make(map[string]*ros_proto.SysStat)

	sub, err := ros_hybrid_go.NewSubscriber(node, "py_test", "sys_msgs/SysInfoStat", &ros_proto.SysInfoStat{}, func(req proto.Message){
		sys := req.(*ros_proto.SysInfoStat)
		info := sys.GetInfo()
		stats[info.GetHostName()] = sys.GetStat()
	})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()


	r := gin.Default()
	r.GET("sys_stat", func(c *gin.Context){
		c.JSON(200, gin.H{
			"stats": stats,
		})
	})

	r.Run(":5151")
}