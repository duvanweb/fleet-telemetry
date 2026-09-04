package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	app := fx.New(Module())

	if err := app.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to start alert-service: %w", err))
	}

	sig := <-app.Wait()
	fmt.Printf("alert-service stopped with code: %v\n", sig.ExitCode)

	if err := app.Stop(ctx); err != nil {
		fmt.Printf("error stopping alert-service: %v\n", err)
	}
}
