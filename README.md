## timewheel-scheduler-demo
## 基于时间轮的任务调度器

### 项目思路与功能
使用时间轮模式实现定时执行任务的功能

1. 支持shell命令类型任务
2. 支持http请求类型任务
3. 支持自定义func类型任务
4. 支持秒级定时执行，并可以设置执行次数或永久执行


```go

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

```
