var builder = WebApplication.CreateSlimBuilder(args);
builder.Logging.ClearProviders();
builder.WebHost.ConfigureKestrel(options => options.AddServerHeader = false);

var app = builder.Build();

app.MapGet("/", () => "Index");

app.Run("http://localhost:5000");
