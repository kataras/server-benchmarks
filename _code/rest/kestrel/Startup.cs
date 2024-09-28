using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Routing;
using System.Collections.Generic;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Http.Features;

namespace netcore
{
    public struct testInput
    {
        public string Name { get; set; }
        public string Language { get; set; }
        public string Id { get; set; }
        public string Bio { get; set; }
        public double Version { get; set; }
    }

    public struct testOutput
    {
        public int id { get; set; }
        public int count { get; set; }
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
            app.Use(async (context, next) =>
            {
                context.Features.Get<IHttpMaxRequestBodySizeFeature>().MaxRequestBodySize = 10 * 1024 *1024 ; // 2MB limit.
                await next.Invoke();
            });
            var routeBuilder = new RouteBuilder(app);
            routeBuilder.MapPost("/{id}", async context =>
            {
                // Follow: https://www.youtube.com/watch?v=gb3zcdZ-y3M
                // https://devblogs.microsoft.com/dotnet/try-the-new-system-text-json-apis/
                // The (new) .NET System.Text.JSON package is
                // x2+ times faster than the Newtonsoft.Json one.
                // So we use that for our bencharks
                // (remember: as fast as possible so we can have real competitors,
                // even if the code does not look as nice as the Iris' one for example). 
                var inputs = await JsonSerializer.DeserializeAsync<List<testInput>>(context.Request.Body);

                var output = new testOutput
                {
                    id = int.Parse(context.GetRouteValue("id").ToString()),
                    count = inputs.Count
                };

                context.Response.Headers.Add("Content-Type", "application/json; charset=utf-8");
                await JsonSerializer.SerializeAsync(context.Response.Body, output);
            });

            var routes = routeBuilder.Build();
            app.UseRouter(routes);
        }
    }
}