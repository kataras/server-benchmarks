using System.Text.Json.Serialization;

var builder = WebApplication.CreateSlimBuilder(args);
builder.Logging.ClearProviders();
builder.WebHost.ConfigureKestrel(options =>
{
    options.AddServerHeader = false;
    options.Limits.MaxRequestBodySize = 2 * 1024 * 1024; // 2MB.
});

var app = builder.Build();

// The {id:int} constraint answers 404 to non-integer ids, matching the
// other frameworks. JSON body binding is case-insensitive by default.
app.MapPost("/{id:int}", (int id, TestInput[] inputs) =>
    Results.Json(new TestOutput(id, inputs.Length, inputs[0].Id)));

app.Run("http://localhost:5000");

sealed record TestInput(string Name, string Language, string Id, string Bio, double Version);

sealed record TestOutput(
    [property: JsonPropertyName("id")] int Id,
    [property: JsonPropertyName("count")] int Count,
    [property: JsonPropertyName("first_id")] string FirstId);
