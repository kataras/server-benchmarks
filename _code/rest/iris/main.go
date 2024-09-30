package main

import "github.com/kataras/iris/v12"

type (
	testInput struct {
		Name     string  `json:"name"`
		Language string  `json:"language"`
		ID       string  `json:"id"`
		Bio      string  `json:"bio"`
		Version  float64 `json:"version"`
	}

	testOutput struct {
		ID      int    `json:"id"`
		Count   int    `json:"count"`
		FirstID string `json:"first_id"`
	}
)

func handler(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(2 * iris.MB) // 2MB

	id := ctx.Params().GetIntDefault("id", 0)

	var in []testInput
	if err := ctx.ReadJSON(&in); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
	})
}

func main() {
	app := iris.New()
	app.Post("/{id:int}", handler)

	app.Listen(":5000")
}
