package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronJob 定义CronJob结构体表示定时任务
type CronJob struct {
	Name     string // 任务名称
	Schedule string // 时间
	Func     func() // 任务功能
}

// CronEngine 定义CronEngine结构体管理定时任务
type CronEngine struct {
	jobs      []*CronJob // 所有的任务
	stop      chan bool  // 是否停止
	reset     chan bool  // 是否重置
	isRunning bool
}

// NewCronEngine 新建一个Cron引擎实例
func NewCronEngine() *CronEngine {
	return &CronEngine{
		jobs:  make([]*CronJob, 0),
		stop:  make(chan bool),
		reset: make(chan bool),
	}
}

// AddJob 添加一个定时任务到Cron引擎
func (ce *CronEngine) AddJob(job *CronJob) {
	ce.jobs = append(ce.jobs, job)
}

// RemoveJob 从Cron引擎中删除一个定时任务
func (ce *CronEngine) RemoveJob(name string) {
	for i, job := range ce.jobs {
		if job.Name == name {
			ce.jobs = append(ce.jobs[:i], ce.jobs[i+1:]...)
			break
		}
	}
}

// Start 启动Cron引擎
// Start 启动Cron引擎
func (ce *CronEngine) Start() {
	if ce.isRunning {
		return
	}
	ce.isRunning = true

	// 使用goroutine在后台运行Cron引擎
	go func() {
		for {
			select {
			// 每分钟触发一次
			case <-time.After(time.Minute):
				for _, job := range ce.jobs {
					if checkSchedule(job.Schedule) {
						job.Func()
					}
				}
				// 接收到停止信号，结束循环并关闭Cron引擎
			case <-ce.stop:
				ce.isRunning = false
				return
				// 接收到重置信号，结束循环并关闭Cron引擎
			case <-ce.reset:
				ce.isRunning = false
				return
			}
		}
	}()
}

// Stop 停止Cron引擎
func (ce *CronEngine) Stop() {
	if !ce.isRunning {
		return
	}
	ce.stop <- true
}

// Reset 重置Cron引擎
func (ce *CronEngine) Reset() {
	if !ce.isRunning {
		return
	}
	ce.reset <- true
}

// 解析Cron表达式，只考虑分钟级别
func checkSchedule(schedule string) bool {
	if schedule == "*" {
		return true
	}

	// 检查是否以 "*/" 开头
	if strings.HasPrefix(schedule, "*/") {
		// 提取 "*/" 后面的数字部分
		numStr := schedule[2:]

		// 检查剩余部分是否是有效的数字
		if num, err := strconv.Atoi(numStr); err == nil {
			currentTime := time.Now()
			minute := currentTime.Minute()
			// 检查当前分钟数是否能被 "*/" 后面的数字整除，能整除代码可以间隔该时间执行
			return minute%num == 0
		}
	} else {
		// 检查是否是有效的数字
		if num, err := strconv.Atoi(schedule); err == nil {
			currentTime := time.Now()
			minute := currentTime.Minute()

			// 检查当前分钟数是否与表达式中的数字相等
			return num == minute
		}
	}

	// 如果表达式不符合上述条件，则返回false
	return false
}

func main() {
	// 示例使用

	// 新建一个Cron引擎实例
	cronEngine := NewCronEngine()

	// 添加定时任务到Cron引擎
	cronEngine.AddJob(&CronJob{
		Name:     "任务1",
		Schedule: "*/1", // 每隔1分钟触发一次
		Func: func() {
			fmt.Println("任务1正在运行", time.Now().Format("2006-01-02 15:04:05"))
		},
	})

	cronEngine.AddJob(&CronJob{
		Name:     "任务2",
		Schedule: "26", // 每小时的第5分钟触发
		Func: func() {
			fmt.Println("任务2正在运行", time.Now().Format("2006-01-02 15:04:05"))
		},
	})

	// 启动Cron引擎
	cronEngine.Start()

	// 等待一段时间后停止Cron引擎
	time.Sleep(10 * time.Minute)
	cronEngine.Stop()
}
