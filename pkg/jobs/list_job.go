package jobs

import (
	"context"
	"os/exec"
	"log"
)

type ListJob struct {
	ctx context.Context
	ResultChan chan interface{}
}

func NewListJob(path string) *ListJob {
	c := context.WithValue(context.TODO(), "path", path)
	return &ListJob{
		ctx: c,
		ResultChan: make(chan interface{}),
	}
}
func (l *ListJob) Do(ctx context.Context) (interface{}, error) {
	log.Printf("processing job.\n")

	path := ctx.Value("path").(string)
	cmdName := "ls"
	cmdArgs := []string{"-al", path}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}

	outPut := string(cmdOut)
	return outPut, nil
}

func (l *ListJob) Context() context.Context {
	return l.ctx
}

func (l *ListJob) Result() chan interface{} {
	return l.ResultChan
}