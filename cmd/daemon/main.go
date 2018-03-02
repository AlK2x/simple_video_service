package main

import (
	"math/rand"
	"os"
	"log"
	"os/signal"
	"syscall"
	"github.com/AlK2x/simple_video_service/packages/model"
	"time"
	"sync"
	"github.com/AlK2x/simple_video_service/packages/ffmpeg"
	"github.com/AlK2x/simple_video_service/packages/repository"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var repo *repository.VideoRepository

type Task struct {
	Id int
	Value *model.VideoItem
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Println("Got SIGINT ...")
	case syscall.SIGTERM:
		log.Println("Got SIGTERM ...")
	}
}

func GenerateTask() *Task {
	if rand.Intn(1000) < 500 {
		return nil
	}

	video, err := repo.GetUnprocessedVideo()
	if err != nil {
		return nil
	}

	return &Task{rand.Intn(10), video}
}

func TaskProvider(stopChan chan struct{}) <-chan *Task {
	tasksChan := make(chan *Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := GenerateTask(); task != nil {
				log.Printf("got the task %v\n", task)
				tasksChan <- task
			} else {
				log.Println("no task for processing, start waiting")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return tasksChan
}

func RunTaskProvider(stopChan chan struct{}) <-chan *Task {
	resultChan := make(chan *Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := TaskProvider(stopTaskProviderChan)
	onStop := func () {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}
	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}

func Worker(tasksChan <-chan *Task, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start handle task %v on worker %v\n", task, name)
		duration, err := ffmpeg.GetVideoDuration(task.Value.Url)
		if err != nil {
			log.Printf("start handle task %v on worker %v\n", task, name)
			continue
		}
		videoKey := task.Value.Item.Id
		thumbnailPath := "content/" + task.Value.Item.Id + "/"
		err = ffmpeg.CreateVideoThumbnail(task.Value.Url, thumbnailPath, 0)
		if err != nil {
			log.Printf("start handle task %v on worker %v\n", task, name)
			continue
		}

		repo.UpdateVideo(videoKey, duration, thumbnailPath)
	}
	log.Printf("stop worker %v\n", name)
}

func RunWorkerPool(stopChan chan struct{}) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := RunTaskProvider(stopChan)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, i)
			wg.Done()
		}(i)
	}
	return &wg
}

func main() {
	db, err := sql.Open("mysql", `root:12345@/simple_video_service`)
	if err != nil {
		log.Fatal(err)
	}
	repo = repository.CreateVideoRepository(db)

	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := getKillSignalChan()
	wg := RunWorkerPool(stopChan)

	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	wg.Wait()
}

