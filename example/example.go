package main


import (
	"context"
	"fmt"
	"github.com/mytools/timewheel-scheduler-demo/pkg"
	"github.com/mytools/timewheel-scheduler-demo/pkg/jobs"
	"time"
)

/*
 参考：https://lk668.github.io/2021/04/05/2021-04-05-%E6%89%8B%E6%8A%8A%E6%89%8B%E6%95%99%E4%BD%A0%E5%A6%82%E4%BD%95%E7%94%A8golang%E5%AE%9E%E7%8E%B0%E4%B8%80%E4%B8%AAtimewheel/
*/

func main() {

	// 初始化时间间隔是1s，一共有60个齿轮的时间轮盘，默认轮盘转动一圈的时间是60s
	tw := timewheel.CreateTimeWheel(1*time.Second, 60, nil)

	// 启动时间轮
	tw.Start()

	// 关闭时间轮
	defer func() {
		tw.Stop()
	}()

	sj := jobs.NewShellJob("kubectl get node", 3)
	fj := jobs.NewFunctionJob(func(ctx context.Context) (interface{}, error) {
		fmt.Println("func jobs...")
		return nil, nil
	}, 2)
	cj, _ := jobs.NewCurlJob("GET", "http://www.baidu.com", "", nil, 2)

	cj1, _ := jobs.NewCurlJob("GET", "http://www.google.com", "", nil, 1)
	if tw.IsRunning() {

		err := tw.AddTask(1*time.Second, "task-shell", time.Now(), sj)
		if err != nil {
			panic(err)
		}

		err = tw.AddTask(2*time.Second, "task-func", time.Now(), fj)
		if err != nil {
			panic(err)
		}

		err = tw.AddTask(1*time.Second, "task-curl", time.Now(), cj)
		if err != nil {
			panic(err)
		}

		err = tw.AddTask(1*time.Second, "task-curl1", time.Now(), cj1)
		if err != nil {
			panic(err)
		}

	} else {
		panic("TimeWheel is not running")
	}

	time.Sleep(3 * time.Second)

	fmt.Println(fj.GetLatestResult())
	fmt.Println(sj.GetLatestResult().Result)
	fmt.Println(cj.GetLatestResult().Result)
	fmt.Println(cj1.GetLatestResult())
	go func() {
		for {
			fmt.Println("in chan")
			fmt.Println(<-sj.GetHistoryResult())
			fmt.Println(<-cj1.GetHistoryResult())
		}
	}()


	<-time.After(time.Second * 60)
	fmt.Println("example test finish...")
}
