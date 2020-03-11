using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace netcore {
    public struct testInput {
        public string email { get; set; }
    }

    public struct testOutput {
        public int id { get; set; }
        public string name { get; set; }
    }

    public class Startup {
        public Startup (IConfiguration configuration) {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices (IServiceCollection services) {
            services.AddRouting ();
        }

        public void Configure (IApplicationBuilder app, IWebHostEnvironment env) {
            var routeBuilder = new RouteBuilder (app);
            routeBuilder.MapPost ("/{id}", context => {
                // Follow: https://www.youtube.com/watch?v=gb3zcdZ-y3M
                // https://devblogs.microsoft.com/dotnet/try-the-new-system-text-json-apis/
                // The (new) .NET System.Text.JSON package is
                // x2+ times faster than the Newtonsoft.Json one.
                // So we use that for our bencharks
                // (remember: as fast as possible se we can have real competitors,
                // even if the code does not look as nice as the Iris' one for example). 
                // var input = JsonSerializer.Deserialize<testInput> (context.Request.Body); <- does not work...
                using (StreamReader stream = new StreamReader (context.Request.Body)) {
                    try {
                        var body = stream.ReadToEndAsync ().Result;
                        var input = JsonSerializer.Deserialize<testInput> (body);

                        var output = new testOutput ();
                        // output.id = (int) context.GetRouteValue ("id"); produces:
                        // Unable to cast object of type 'System.String' to type 'System.Int32'.
                        // Another disadvantage, in Iris you could get the value as defined in routing,
                        // e.g. {id:int} without convert it again and again.... Anyway
                        output.id = int.Parse (context.GetRouteValue ("id").ToString ());
                        output.name = input.email;

                        context.Response.Headers.Add ("Content-Type", "application/json; charset=utf-8");
                        // default json options: minified output.
                        return context.Response.WriteAsync (JsonSerializer.Serialize (output));
                    } catch (Exception e) {
                        return context.Response.WriteAsync (e.Message);
                    }
                }
            });
            var routes = routeBuilder.Build ();
            app.UseRouter (routes);
        }
    }
}