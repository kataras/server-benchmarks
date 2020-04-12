using System;
using System.Text.Json;
using System.Threading.Tasks;
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
        private static JsonSerializerOptions JsonSerializerOptions = new JsonSerializerOptions();

        public IConfiguration Configuration { get; }

        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddRouting();
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            var routeBuilder = new RouteBuilder(app);
            routeBuilder.MapPost("/{id}", context =>
            {
                Span<byte> reqBytes = new byte[(int)context.Request.ContentLength.GetValueOrDefault(32)];

                var reqLen = context.Request.Body.Read(reqBytes);
                // Follow: https://www.youtube.com/watch?v=gb3zcdZ-y3M
                // https://devblogs.microsoft.com/dotnet/try-the-new-system-text-json-apis/
                // The (new) .NET System.Text.JSON package is
                // x2+ times faster than the Newtonsoft.Json one.
                // So we use that for our bencharks
                // (remember: as fast as possible so we can have real competitors,
                // even if the code does not look as nice as the Iris' one for example). 
                var input = JsonSerializer.Deserialize<testInput>(reqBytes.Slice(0, reqLen), JsonSerializerOptions);

                var output = new testOutput
                {
                    // output.id = (int) context.GetRouteValue ("id"); produces:
                    // Unable to cast object of type 'System.String' to type 'System.Int32'.
                    // Another disadvantage, in Iris you could get the value as defined in routing,
                    // e.g. {id:int} without convert it again and again.... Anyway
                    id = int.TryParse(context.GetRouteValue("id").ToString(), out int t) ? t : 0,
                    name = input.email
                };

                context.Response.Headers.Add("Content-Type", "application/json; charset=utf-8");

                // default json options: minified output.
                ReadOnlySpan<byte> outputBytes = JsonSerializer.SerializeToUtf8Bytes(output, typeof(testOutput), JsonSerializerOptions);

                context.Response.Body.Write(outputBytes);

                return Task.CompletedTask;
            });
            var routes = routeBuilder.Build();
            app.UseRouter(routes);
        }
    }
}