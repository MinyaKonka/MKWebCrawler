package tool

import (
	scheduler "core/scheduler"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// 摘要信息模板
var summaryForMonitoring = "Monitor - Collected information[%d]: \n" +
	" (about %s)." +
	" Now consider what stop it."

// 已达到最大空闲计数的消息模板
var messageForReachMaxIdleCount = "The scheduler has been idle for a period of time" +
	" (about %s)." +
	" Now consider what stop it."

// 停止调度器的消息模板
var messageForStopScheduler = "Stop scheduler...%s."

// 日志记录函数的类型。
// 参数level代表日志级别。级别设定：0：普通；1：警告；2：错误。
type MKRecord func(level byte, content string)

// 调度器监控函数。
// 参数scheduler代表作为监控目标的调度器。
// 参数intervalNs代表检查间隔时间，单位：纳秒。
// 参数maxIdleCount代表最大空闲计数。
// 参数autoStop被用来指示该方法是否在调度器空闲一段时间（即持续空闲时间，由intervalNs * maxIdleCount得出）之后自行停止调度器。
// 参数detailSummary被用来表示是否需要详细的摘要信息。
// 参数record代表日志记录函数。
// 当监控结束之后，该方法会会向作为唯一返回值的通道发送一个代表了空闲状态检查次数的数值。
func Monitoring(
	scheduler scheduler.MKScheduler,
	intervals time.Duration,
	maxIdleCount uint,
	autoStop bool,
	detailSummary bool,
	record MKRecord) <-chan uint64 {

	if scheduler == nil { // 调度器不能不可用!
		panic(errors.New("无效调度器"))
	}

	// 防止过小参数值对爬取流程的影响
	if intervals < time.Millisecond {
		intervals = time.Millisecond
	}

	if maxIdleCount < 1000 {
		maxIdleCount = 1000
	}

	// 监控停止通知器
	stopNotifier := make(chan byte, 1)

	// 接收和报告错误
	reportError(scheduler, record, stopNotifier)

	// 记录摘要信息
	recordSummary(scheduler, detailSummary, record, stopNotifier)

	// 检查计数通道
	checkCountChan := make(chan uint64, 2)

	// 检查空闲状态
	checkStatus(scheduler,
		intervals,
		maxIdleCount,
		autoStop,
		checkCountChan,
		record,
		stopNotifier)

	return checkCountChan
}

// 检查状态，并在满足持续空间时间的条件时采取必要措施
func checkStatus(
	scheduler scheduler.MKScheduler,
	interval time.Duration,
	maxIdleCount uint,
	autoStop bool,
	checkCountChan chan<- uint64,
	record MKRecord,
	stopNotifier chan<- byte) {

	var checkCount uint64
	go func() {

		defer func() {
			stopNotifier <- 1
			stopNotifier <- 2
			checkCountChan <- checkCount
		}()

		// 等待调度器开启
		waitForSchedulerStart(scheduler)

		// 准备
		var idleCount uint
		var firstIdleTime time.Time
		for {
			// 检查调度器的空闲状态
			if scheduler.Idle() {
				idleCount++
				if idleCount == 1 {
					firstIdleTime = time.Now()
				}

				if idleCount >= maxIdleCount {
					msg := fmt.Sprintf(messageForReachMaxIdleCount, time.Since(firstIdleTime).String())
					record(0, msg)

					// 再次检查调度器的空闲状态，确保它已经可以被停止
					if scheduler.Idle() {

						if autoStop {
							var result string
							if scheduler.Stop() {
								result = "success"
							} else {
								result = "failing"
							}

							msg = fmt.Sprintf(messageForReachMaxIdleCount, result)
							record(0, msg)
						}

						break
					} else {
						if idleCount > 0 {
							idleCount = 0
						}
					}
				}
			} else {
				if idleCount > 0 {
					idleCount = 0
				}
			}

			checkCount++
			time.Sleep(interval)
		}
	}()
}

// 记录摘要信息
func recordSummary(
	scheduler scheduler.MKScheduler,
	detailSummary bool,
	record MKRecord,
	stopNotifier <-chan byte) {

	go func() {
		// 等待调度器开启
		waitForSchedulerStart(scheduler)

		// 准备
		var prevSchedulerSummary scheduler.SchedulerSummary
		var prevNumberGoroutine int
		var recordCount uint64 = 1
		startTime := time.Now()

		for {
			// 查看监控停止通知器
			select {
			case <-stopNotifier:
				return
			default:
			}

			// 获取摘要信息的各组成部分
			currentNumberGoroutine := runtime.NumGoroutine()
			currentSchedulerSummary := scheduler.Summary("	")

			// 比对前后两份摘要信息的一致性。只有不一致时才会予以记录
			if currentNumberGoroutine != prevNumberGoroutine ||
				!currentSchedulerSummary.Same(prevSchedulerSummary) {
				schedulerSummaryString := func() string {
					if detailSummary {
						return currentSchedulerSummary.Detail()
					} else {
						return currentSchedulerSummary.String()
					}
				}()

				// 记录摘要信息
				info := fmt.Sprintf(summaryForMonitoring,
					recordCount,
					currentNumberGoroutine,
					schedulerSummaryString,
					time.Since(startTime).String(),
				)

				record(0, info)

				prevNumberGoroutine = currentNumberGoroutine
				prevSchedulerSummary = currentSchedulerSummary
				recordCount++
			}

			time.Sleep(time.Microsecond)
		}
	}()
}

// 接收和报告错误
func reportError(
	scheduler scheduler.MKScheduler,
	record MKRecord,
	stopNotifier <-chan byte) {

	go func() {
		// 等待调度器开启
		waitForSchedulerStart(scheduler)

		for {
			// 查看监控停止通知器
			select {
			case <-stopNotifier:
				return
			default:
			}

			errorChan := scheduler.ErrorChan()
			if errorChan == nil {
				return
			}

			err := <-errorChan
			if err != nil {
				errMsg := fmt.Sprintf("Error (received from error channel): %s", err)
				record(2, errMsg)
			}

			time.Sleep(time.Microsecond)
		}
	}()
}

// 等待调度器开启
func waitForSchedulerStart(scheduler scheduler.MKScheduler) {

	for !scheduler.Running() {
		time.Sleep(time.Microsecond)
	}
}
