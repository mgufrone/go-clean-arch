package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"os"
	"recipes/app"
	"recipes/app/common"
	"recipes/app/config"
	"recipes/domains/recipe"
	"recipes/worker/worker"
	"time"
)

func initiateWorker(handler recipe.IRecipeUseCase) worker.RecipeWorker {
	return worker.NewRecipeWorker(context.Background(), handler)
}
type WorkerContext struct {

}
func runWorker(lc fx.Lifecycle, wk worker.RecipeWorker, cfg *config.Config)  {
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle: 5,
		Wait: true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
		},
	}
	workerPool := work.NewWorkerPool(WorkerContext{}, 5, string(common.QueueDefault), redisPool)
	workerPool.Middleware(wk.Logger)
	workerPool.Job(string(common.ProcessLike), wk.ProcessLike)
	workerPool.Job(string(common.ProcessDislike), wk.ProcessDislike)
	workerPool.Job(string(common.ProcessView), wk.ProcessView)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			workerPool.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if workerPool == nil {
				panic(errors.New("worker pool is not initiated"))
			}
			workerPool.Stop()
			return nil
		},
	})
}
func main() {
	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}
	a := app.App(
		fx.Provide(initiateWorker),
		fx.Invoke(runWorker),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	if err := a.Start(ctx); err != nil {
		panic(err)
	}
	<- a.Done()
}
