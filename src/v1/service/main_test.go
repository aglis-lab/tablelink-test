package service_test

import (
	"context"
	"os"
	"tablelink/src/app"
	"testing"
)

func TestMain(m *testing.M) {
	os.Chdir("../../../")

	app.Init(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)
}
