using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace netcore
{
    public struct testInput
    {
        public string email { get; set; }
    }

    public struct testOutput
    {
        public int id { get; set; }
        public string name { get; set; }
    }

    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddRouting();
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            var routeBuilder = new RouteBuilder(app);
            routeBuilder.MapPost("/{id}", context =>
            {
                // Follow: https://www.youtube.com/watch?v=gb3zcdZ-y3M
                // https://devblogs.microsoft.com/dotnet/try-the-new-system-text-json-apis/
                // The (new) .NET System.Text.JSON package is
                // x2+ times faster than the Newtonsoft.Json one.
                // So we use that for our bencharks
                // (remember: as fast as possible so we can have real competitors,
                // even if the code does not look as nice as the Iris' one for example). 
                var input = JsonSerializer.DeserializeAsync<testInput>(context.Request.Body).Result;

                var output = new testOutput
                {
                    // output.id = (int) context.GetRouteValue ("id"); produces:
                    // Unable to cast object of type 'System.String' to type 'System.Int32'.
                    // Another disadvantage, in Iris you could get the value as defined in routing,
                    // e.g. {id:int} without convert it again and again.... Anyway
                    id = int.Parse(context.GetRouteValue("id").ToString()),
                    name = input.email
                };

                context.Response.Headers.Add("Content-Type", "application/json; charset=utf-8");
                // default json options: minified output.
                return JsonSerializer.SerializeAsync(context.Response.Body, output);
            });
            var routes = routeBuilder.Build();
            app.UseRouter(routes);
        }
    }
}