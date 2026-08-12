var builder = WebApplication.CreateSlimBuilder(args);
builder.Logging.ClearProviders();
builder.WebHost.ConfigureKestrel(options => options.AddServerHeader = false);

var app = builder.Build();

app.MapGet("/hello/{name}", (string name) => $"Hello {name}");

app.Run("http://localhost:5000");
