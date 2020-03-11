package main

import "github.com/kataras/iris/v12"

type (
	testInput struct {
		Email string `json:"email"`
	}

	testOutput struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func handler(ctx iris.Context) {
	id := ctx.Params().GetIntDefault("id", 0)

	var in testInput
	if err := ctx.ReadJSON(&in); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(testOutput{
		ID:   id,
		Name: in.Email,
	})
}

func main() {
	app := iris.New()
	app.Handle(iris.MethodPost, "/{id:int}", handler)
	app.Listen(":5000", iris.WithOptimizations)
}
