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
		panic(fmt.Errorf("failed to start vehicle-service: %w", err))
	}

	sig := <-app.Wait()
	fmt.Printf("vehicle-service stopped with code: %v\n", sig.ExitCode)

	if err := app.Stop(ctx); err != nil {
		fmt.Printf("error stopping vehicle-service: %v\n", err)
	}
}
