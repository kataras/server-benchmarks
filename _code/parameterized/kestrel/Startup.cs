using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace netcore {
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
            routeBuilder.MapGet ("/hello/{name}", context => {
                return context.Response.WriteAsync ("Hello " + context.GetRouteValue ("name"));
            });
            var routes = routeBuilder.Build ();
            app.UseRouter (routes);
        }
    }
}